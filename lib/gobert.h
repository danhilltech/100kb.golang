#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct SharedSentenceEmbeddingModel SharedSentenceEmbeddingModel;
typedef struct SharedKeywordExtractionModel SharedKeywordExtractionModel;
typedef struct SharedZeroShotModel SharedZeroShotModel;

struct SharedSentenceEmbeddingModel *new_sentence_embedding(void);
void drop_sentence_embedding(struct SharedSentenceEmbeddingModel *ptr);

uint8_t *sentence_embedding(struct SharedSentenceEmbeddingModel *ptr,
                            const unsigned char *req,
                            size_t *req_size,
                            size_t *out_size);

struct SharedKeywordExtractionModel *new_keyword_extraction(void);
void drop_keyword_extraction(struct SharedKeywordExtractionModel *ptr);

uint8_t *keyword_extraction(struct SharedKeywordExtractionModel *ptr,
                            const unsigned char *req,
                            size_t *req_size,
                            size_t *out_size);

struct SharedZeroShotModel *new_zero_shot(void);
void drop_zero_shot(struct SharedZeroShotModel *ptr);

uint8_t *zero_shot(struct SharedZeroShotModel *ptr,
                   const unsigned char *req,
                   size_t *req_size,
                   size_t *out_size);

void drop_bytesarray(uint8_t *ptr);