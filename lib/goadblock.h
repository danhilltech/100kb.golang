#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct AdblockEngine AdblockEngine;

struct AdblockEngine *new_adblock(const unsigned char *req,
                                  size_t *req_size);
void drop_adblock(struct AdblockEngine *ptr);

uint8_t *filter(struct AdblockEngine *ptr,
                const unsigned char *req,
                size_t *req_size,
                size_t *out_size);

void drop_bytesarray(uint8_t *ptr);