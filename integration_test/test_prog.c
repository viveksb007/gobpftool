// Minimal eBPF program for integration testing
// Use kernel headers compatible definitions

typedef unsigned int __u32;
typedef unsigned long long __u64;

// BPF map type
#define BPF_MAP_TYPE_ARRAY 2

// XDP return codes
#define XDP_PASS 2

// Helper function
static void *(*bpf_map_lookup_elem)(void *map, const void *key) = (void *)1;

// Map definition using BTF-style
struct {
    int (*type)[BPF_MAP_TYPE_ARRAY];
    __u32 *key;
    __u64 *value;
    int (*max_entries)[1];
} test_map __attribute__((section(".maps")));

struct xdp_md {
    __u32 data;
    __u32 data_end;
};

// Simple XDP program that passes all packets
__attribute__((section("xdp"), used))
int test_prog(struct xdp_md *ctx) {
    __u32 key = 0;
    __u64 *value = bpf_map_lookup_elem(&test_map, &key);
    if (value) {
        __sync_fetch_and_add(value, 1);
    }
    return XDP_PASS;
}

char _license[] __attribute__((section("license"), used)) = "GPL";
