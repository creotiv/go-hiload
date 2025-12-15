# **Why Buffered Writes Are Dangerous in High-Load Systems**

When engineers first benchmark file I/O, buffered writes look fantastic.
They show huge throughput, extremely low latency, and nearly zero CPU overhead.

But this performance is an illusion.

For high-load systems — especially those using **WAL**, **log segments**, **append-only storage**, or **durability guarantees** — buffered write can quietly destroy your reliability and ruin system stability at scale.

This article explains **why buffered I/O becomes dangerous** and **when (rarely) you can still use it**.

---

# **1. Buffered Write Loses Data on Crash**

Buffered I/O writes into RAM (page cache), not the disk.
The OS decides *when* to flush the pages.

If your process, VM, or hardware dies before flush — all buffered data disappears.

For logs/WAL this is catastrophic:
**if a commit is acknowledged but not on disk, it is lost forever.**

That’s why real high-reliability systems avoid buffered I/O:

* PostgreSQL WAL → **O_DIRECT** or **fdatasync on each commit**
* Kafka commit log → **O_DIRECT**
* ClickHouse parts → **O_DIRECT**
* VictoriaLogs → **O_DIRECT**
* Prometheus WAL → `write + fdatasync`
* Loki → linear WAL segments + controlled fsync

**Rule #1 of logging:**

> *“A log entry is considered written only when it physically reaches disk.”*

Buffered writes simply cannot guarantee this.

---

# **2. Buffered I/O Causes Write Amplification & Fragmentation**

Buffered writes force the kernel to do a lot of work unrelated to your actual write:

* read-modify-write of partially dirty pages
* maintenance of dirty page lists
* periodic flush timers
* page-level fragmentation
* inode churn
* write amplification inside the filesystem
* extra work for the NVMe controller

Direct I/O does none of this. It writes **exact bytes to exact Local Block Addresses**.

Buffered I/O always costs more in the long run — more CPU, more disk ops, more latency spikes.

---

# **3. Page Cache Is a CPU Resource — Logs Should Not Pollute It**

Logs and WAL typically have this pattern:

* **write a lot**
* **read almost never**

Putting logs into page cache is wasteful and harmful:

* they evict hot code/data pages of your application
* the kernel must manage large dirty page sets
* more cache churn increases CPU stalls
* JVM/Go runtimes suffer increased GC pauses
* overall latency becomes unpredictable

Direct I/O leaves cache space for user-mode processes, not logs.

---

# **4. Buffered Write Under High Load Triggers Write Storms**

Under light load buffered writes are “fine”.
But under **200K–500K RPS** they become unstable.

The kernel starts to:

* batch and flush pages unevenly
* execute massive writeback bursts
* stall all writers during flush
* cause dramatic tail-latency jitter
* lock filesystem metadata
* temporarily freeze the process

This leads to unpredictable **99th/99.9th percentile latency**, which is fatal for:

* ingestion systems
* real-time pipelines
* distributed storage nodes
* HA systems with consensus
* WAL committers

Direct I/O ensures **stable tail latency** because it avoids all writeback machinery.

---

# **5. Buffered Writes Hide the True Cost of I/O**

Buffered I/O gives engineers the wrong idea:

* “writes are cheap”
* “the disk is fast”
* “latency is stable”

Then suddenly:

* writeback daemon runs at 3AM and stalls the whole node
* `fsync()` blocks for seconds
* latency jumps ×1000
* NVMe hits thermal throttle and everything collapses

Buffered I/O **masks the real cost**.
Direct I/O **shows it immediately**, allowing real capacity planning.

---

# **But wait, if we write directly to disk, why we still need FSYNC?**

O_DIRECT does NOT guarantee durability. This is the key point that many engineers misunderstand.

O_DIRECT bypasses the page cache, but it does NOT bypass:

* filesystem journaling
* filesystem metadata
* NVMe write buffers
* controller caches
* drive-level DRAM
* reordering by kernel I/O scheduler
* reordering inside the NVMe firmware

Therefore:

🚨 O_DIRECT != durable write

🚨 O_DIRECT != write is on disk

🚨 O_DIRECT != crash-safe

O_DIRECT only means:

“Write directly to the block device without using page cache.”

It does NOT mean:

“Write is persisted on stable storage.”

That’s why you still need fsync / fdatasync.

# Why O_DIRECT still needs fsync / fdatasync?

## 1. Filesystem metadata is still buffered

Even with O_DIRECT, the filesystem metadata (inode timestamps, allocation bitmaps, log extents, directory updates, journal blocks) must still be flushed.

Example:
```
O_DIRECT write  → writes your data
BUT
fs still updates:
- inode
- file size
- allocation
- journal (for ext4, xfs)
```

Those updates sit in buffered metadata and will NOT persist without fsync.

## 2. NVMe drives reorder writes unless you fsync

NVMe controllers buffer and reorder writes for performance.

They commit writes to internal DRAM, not NAND, unless:

* flush command is issued (via fsync)
* FUA bit set (not used by Go’s Pwrite)
* write barrier is issued

Without fsync:
```
power loss → NVMe loses DRAM → DATA corrupted
```

This is why:

* PostgreSQL uses O_DIRECT + fdatasync
* Kafka uses O_DIRECT + fsync policies
* ClickHouse uses O_DIRECT + fsync-on-part-finalization
* Redis AOF uses fsync policy (always/everysec/no)

## 3. O_DIRECT avoids page cache — NOT ordering

Linux still does reordering internally unless a flush barrier is issued.

fsync and fdatasync issue the barrier.

## 4. File contents may be on disk — but NOT guaranteed to be readable after crash

ext4 and xfs use journaling:

* journal writes can be delayed
* commit records not flushed
* data + metadata may be inconsistent

Only fsync ensures:
```
data → disk
metadata → disk
journal commit → disk
```

## Aligment

Data alignment is required when using O_DIRECT because direct I/O operations involve Direct Memory Access (DMA) between the user-space buffer and the storage device, bypassing the kernel's page cache. 

The key reasons for this requirement are:

* **Hardware Compatibility**: Storage devices perform I/O in fixed-size units called sectors (or logical blocks), which are typically 512 bytes or 4096 bytes (4KB). The hardware disk controller expects data transfers to start and end on these specific boundaries.

* **Bypassing the Kernel**: In normal, buffered I/O, the kernel handles the "fixing up" of unaligned or partial requests by using an intermediate cache (page cache) and performing read-modify-write operations if necessary. With O_DIRECT, the application takes full responsibility for coordinating the I/O, and there is no kernel buffer to align the data automatically.
* **Zero-Copy Efficiency**: The primary benefit of O_DIRECT is eliminating data copies between the user buffer and kernel memory. This "zero-copy" approach requires that the user's memory buffer addresses align with the physical memory addresses that the DMA controller needs to access the storage device directly.

* **Performance and Simplicity**: Unaligned transfers would require the kernel to implement bounce buffers and perform read-modify-write cycles, which would add overhead and negate the performance benefits of using O_DIRECT in the first place. 


## Practical example: why O_DIRECT without fsync is dangerous

Imagine WAL/log chunk:
```
seg_00001.wal
```

You write 64 KB aligned blocks using O_DIRECT.

Crash happens.

You might get:

* a file that looks correct but has incomplete metadata
* block that was reordered after later blocks (NVMe)
* extents allocated but not updated in allocation tree
* inode not pointing to correct blocks
* corruption in the middle of the segment

That WAL is **unrecoverable.**


# **When Buffered Writes *Can* Be Used in High-Load Systems**

Buffered writes are not evil — just dangerous for durability-critical paths.

Use buffered writes only when **data loss is acceptable** and **latency spikes do not break correctness**.

✔ Suitable cases:

* debug logs
* audit logs **without durability requirements**
* analytics files written asynchronously
* temporary data that can be recomputed
* cold-storage batch processing
* low-priority event pipeline for optional metrics

Buffered write makes sense only when:

* losing the last seconds/minutes of data is fine
* the consumer tolerates jitter
* recovery after crash is trivial
* write ordering does not matter

Think: monitoring dashboards, transient caches, local CSV exports, etc.

---

# **When Buffered Writes Must NEVER Be Used**

**NEVER use buffered write for:**

* WAL logs in a database or queue
* ingestion logs (Loki, VictoriaLogs, Prometheus, ClickHouse)
* replication logs
* consensus logs (Raft/Paxos)
* anything that is later uploaded to S3
* append-only segments with strict ordering
* recovery-critical metadata
* anything requiring deterministic durability guarantees

If your system needs:

* crash safety
* predictable tail latency
* linear write performance
* stable NVMe load
* zero writeback surprises
* minimal filesystem overhead
* correct ordering guarantees

→ **Use O_DIRECT or an equivalent direct-write API.**

Period.

---

# **Summary: The Rule of Thumb**

Buffered writes are fast only because the OS lies to you.
Under high load, this lie becomes dangerous.

Use **direct/O_DIRECT-style writes** when building:

* WAL
* ingestion pipelines
* log segments
* replication logs
* distributed storage
* anything recovery-critical

Use **buffered writes** only when data loss is OK and jitter does not break correctness.


# What results would you see on a real physical Linux machine?
### On physical x86 Linux with NVMe:
```
BufferedWrite     ~200 ns/op
DirectWrite       ~600–900 ns/op
```

### On SSD SATA:
```
BufferedWrite     ~200 ns/op
DirectWrite       ~4,000–10,000 ns/op
```

### On HDD:
```
BufferedWrite     ~200 ns/op
DirectWrite       ~20 ms (!!!)
```

**Inside a VM (Lima) results ALWAYS inflate direct I/O latency massively.
That’s normal behavior.**