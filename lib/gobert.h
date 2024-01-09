#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct Model Model;

struct Model *new_sentence_embedding(void);
void drop_sentence_embedding(struct Model *ptr);

void sentence_embedding(struct Model *ptr, const char *const *strs, void *dst);
