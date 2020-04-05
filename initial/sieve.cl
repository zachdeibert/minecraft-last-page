typedef struct {
    uint prefix_len;
    uint base_len;
    uint suffix_len;
} config_t;

typedef struct {
    uint pos;
    uint size;
    uchar buf[0];
} output_list_t;

__constant const uchar salt[] = ":why_so_salty#LazyCrypto";
__constant const uint salt_len = 24;
__constant const uchar alphabet[] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()";
__constant const uint alphabet_len = 72;
__constant const uint target_hash = 0x4A86C76A;
__constant const uint target_mask = 0xFFFFFF7F;

__constant const uint empty_hash[] = {
    0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19
};
__constant const uint sha_k[] = {
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
    0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
    0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
    0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
    0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
    0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2
};

inline uint rotr32(uint val, uint count) {
    return rotate(val, 32 - count);
}

void sha256_block(uint *hash, uint *chunk) {
    uint w[64];
    for (int i = 0; i < 16; ++i) {
        w[i] = chunk[i];
    }
    for (int i = 16; i < 64; ++i) {
        uint s0 = rotr32(w[i - 15], 7) ^ rotr32(w[i - 15], 18) ^ (w[i - 15] >> 3);
        uint s1 = rotr32(w[i - 2], 17) ^ rotr32(w[i - 2], 19) ^ (w[i - 2] >> 10);
        w[i] = w[i - 16] + s0 + w[i - 7] + s1;
    }
    uint a = hash[0];
    uint b = hash[1];
    uint c = hash[2];
    uint d = hash[3];
    uint e = hash[4];
    uint f = hash[5];
    uint g = hash[6];
    uint h = hash[7];
    for (int i = 0; i < 64; ++i) {
        uint s1 = rotr32(e, 6) ^ rotr32(e, 11) ^ rotr32(e, 25);
        uint ch = (e & f) ^ ((~e) & g);
        uint temp1 = h + s1 + ch + sha_k[i] + w[i];
        uint s0 = rotr32(a, 2) ^ rotr32(a, 13) ^ rotr32(a, 22);
        uint maj = (a & b) ^ (a & c) ^ (b & c);
        uint temp2 = s0 + maj;
        h = g;
        g = f;
        f = e;
        e = d + temp1;
        d = c;
        c = b;
        b = a;
        a = temp1 + temp2;
    }
    hash[0] += a;
    hash[1] += b;
    hash[2] += c;
    hash[3] += d;
    hash[4] += e;
    hash[5] += f;
    hash[6] += g;
    hash[7] += h;
}

void sha256_single(uint *hash, uchar *data, uint data_len) {
    if (data_len > 55) {
        data_len = 55;
        printf("Error: data_len (%d) > 55\n", data_len);
    }
    for (int i = 0; i < 8; ++i) {
        hash[i] = empty_hash[i];
    }
    uint chunk[16];
    for (int i = 0; i < 16; ++i) {
        chunk[i] = 0;
    }
    uchar *chunk_bytes = (uchar *) chunk;
    uchar *chunk_it = chunk_bytes;
    for (int i = 0; i < data_len; ++i) {
        *chunk_it++ = data[i];
    }
    *chunk_it = 0x80;
    chunk_bytes[62] = (data_len >> 5) & 0xFF;
    chunk_bytes[63] = (data_len << 3) & 0xFF;
    for (int i = 0; i < 16; ++i) {
        chunk[i] = ((chunk[i] & 0x000000FF) << 24) |
                   ((chunk[i] & 0x0000FF00) <<  8) |
                   ((chunk[i] & 0x00FF0000) >>  8) |
                   ((chunk[i] & 0xFF000000) >> 24);
    }
    sha256_block(hash, chunk);
}

bool check_hash(uchar *str, uint str_len) {
    uint hash[8];
    sha256_single(hash, str, str_len);
    return (hash[0] & target_mask) == (target_hash & target_mask);
}

void report_result(__global volatile output_list_t *outs, uchar *str, uint str_len) {
    uint start = atomic_add(&outs->pos, str_len);
    if (start + str_len >= outs->size) {
        printf("Error: output buffer overflow\n");
    }
    __global volatile uchar *buf = outs->buf + start;
    for (int i = 0; i < str_len; ++i) {
        *buf++ = *str++;
    }
}

__kernel void last_page_sieve(__global config_t *cfg, __global uchar *prefix, __global volatile output_list_t *outs) {
    int id = get_global_id(0);
    if (id == 0) {
        outs->pos = 0;
    }
    barrier(CLK_GLOBAL_MEM_FENCE);
    const uint str_len = cfg->prefix_len + cfg->base_len + cfg->suffix_len;
    if (str_len > 55 - salt_len) {
        printf("Error: string too long\n");
        return;
    }
    uchar str[55];
    for (int i = 0; i < cfg->prefix_len; ++i) {
        str[i] = prefix[i];
    }
    uchar *base = str + cfg->prefix_len;
    uchar *suffix = base + cfg->base_len;
    for (int i = 0; i < cfg->suffix_len; ++i) {
        *suffix++ = alphabet[id % alphabet_len];
        id /= alphabet_len;
    }
    for (int i = 0; i < salt_len; ++i) {
        *suffix++ = salt[i];
    }
    id = get_global_id(0);
    int idxs[55];
    for (int i = 0; i < cfg->base_len - 2; ++i) {
        idxs[i] = 0;
        base[i] = alphabet[0];
    }
    idxs[cfg->base_len - 1] = -1;
    while (1) {
        for (int i = cfg->base_len - 1; i >= 0; --i) {
            if (++idxs[i] < alphabet_len) {
                base[i] = alphabet[idxs[i]];
                if (check_hash(str, str_len + salt_len)) {
                    report_result(outs, str, str_len);
                }
                break;
            } else if (i == 0) {
                return;
            } else {
                idxs[i] = 0;
                base[i] = alphabet[0];
            }
        }
    }
}
