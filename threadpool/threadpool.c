#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<pthread.h>
#include<string.h>
#include<stdbool.h>

// task node..
typedef struct task_node {
    void (* func) (void* args);
    void *args;
} task_node;

typedef struct {
   task_node* q;
   int cap;
   int len;
   int readAt, insertAt;

} queue;

bool init_queue(queue* q, int cap) {
    if (cap <= 0 || cap >= (int) 1e6) { // for now, no long capacity
        cap = 1;
    }
    q->cap = cap;
    q->len = 0;
    q->readAt = q->insertAt = 0;
    task_node *qq = (task_node*) malloc(sizeof(task_node) * q->cap);
    if (qq == NULL) {
        return false;
    }
    q->q = qq;
    return true;
}

bool isFull(queue* q) {
    return q->len == q->cap;
}

bool isEmpty(queue* q) {
    return q->len == 0;
}

bool resize(queue* q) {
    if (!isFull(q)) {
        return false;
    }
    int newCap = 2 * q->cap;
    if (newCap <= 0) {   // greater than 1e9
        return false;
    }
    task_node* newQ = (task_node*) malloc(sizeof(task_node) * newCap);
    if (newQ == NULL) {
        return false;
    }
    for (int i = 0; i < q->len; i++) {
        newQ[i] = q->q[(q->readAt + i) % q->cap];
    }
    free(q->q);
    q->q = newQ;
    q->cap = newCap;
    q->insertAt = q->len;
    q->readAt = 0;
    return true;
}

bool push(queue* q, task_node val) {
    if (isFull(q)) {
        bool ok = resize(q);
        if (!ok) return ok;
    }
    q->q[q->insertAt] = val;
    q->insertAt = (q->insertAt + 1) % q->cap; 
    q->len++;
    return true;
}

bool pop(queue* q, task_node* val) {
    if(isEmpty(q)) {
        return false;
    }
    *val = q->q[q->readAt];
    q->readAt = (q->readAt + 1) % q->cap;
    q->len--;
    return true;
}

void free_queue(queue* q) {
    if(q == NULL) return;
    if(q->q != NULL) free(q->q);
    if(q != NULL) free(q);
}

const int MAX_THREAD_COUNT = 4;
typedef struct {
    pthread_t thread[MAX_THREAD_COUNT];
    pthread_mutex_t lock;
    pthread_cond_t cond;
    queue* q;
    bool shutdown;
} thread_pool;

void* start_pool(void* tPool);
void destroy_pool(thread_pool* pool, int threadCount);

thread_pool* create_pool(int queueSize) {
    thread_pool* pool = (thread_pool*) malloc(sizeof(thread_pool));
    if (pool == NULL) {
        return NULL;
    }
    pool->shutdown = false;
    pool->q = (queue*) malloc(sizeof(queue));
    if (pool->q == NULL) {
        free(pool);
        return NULL;
    }
    if (!init_queue(pool->q, queueSize)) {
        free_queue(pool->q);
        free(pool);
        return NULL;
    }
    if (pthread_mutex_init(&pool->lock, NULL) != 0) {
        free_queue(pool->q);
        free(pool);
        return NULL;
    }
    if (pthread_cond_init(&pool->cond, NULL) != 0) {
        pthread_mutex_destroy(&pool->lock);
        free_queue(pool->q);
        free(pool);
        return NULL;
    }
    for (int i = 0; i < MAX_THREAD_COUNT; i++) {
        if (pthread_create(&(pool->thread[i]), NULL, start_pool, (void*) pool) != 0) {
            destroy_pool(pool, i);
            return NULL;
        }
    }
    return pool;
}

bool push_task(thread_pool* pool, void (* func) (void* args), void *args) {
    if (pool == NULL || func == NULL) {
        return false;
    }
    pthread_mutex_lock(&pool->lock);
    if (pool->shutdown) {
        pthread_mutex_unlock(&pool->lock);
        return false;
    }
    task_node node = {
        .func = func,
        .args = args
    };
    if (!push(pool->q, node)) {
        pthread_mutex_unlock(&pool->lock);
        return false;
    }
    pthread_cond_signal(&pool->cond);
    pthread_mutex_unlock(&pool->lock);
    return true;
}

void destroy_pool(thread_pool* pool, int threadCount) {
    pthread_mutex_lock(&pool->lock);
    pool->shutdown = true;
    pthread_cond_broadcast(&pool->cond);
    pthread_mutex_unlock(&pool->lock);
    for (int i = 0; i < threadCount; i++) {
        pthread_join(pool->thread[i], NULL);
    }
    pthread_mutex_destroy(&pool->lock);
    pthread_cond_destroy(&pool->cond);
    free_queue(pool->q);
    free(pool);
}

void* start_pool(void* tPool) {
    thread_pool* pool = (thread_pool*) (tPool);
    task_node task;
    while(1) {
        pthread_mutex_lock(&pool->lock);
        while(isEmpty(pool->q) && !pool->shutdown) {
            pthread_cond_wait(&pool->cond, &pool->lock);
        }
        if (pool->shutdown && isEmpty(pool->q)) {
            pthread_mutex_unlock(&pool->lock);
            break;
        }
        pop(pool->q, &task);
        pthread_mutex_unlock(&pool->lock);
        // executing task...
        (*(task.func))(task.args);
    }
    pthread_exit(NULL);
    return NULL;
}

void print_task(void* arg) {
    int id = *(int*)arg;
    printf("[Thread %lu] Task #%d starting...\n", (unsigned long)pthread_self(), id);
    sleep(1);
    printf("[Thread %lu] Task #%d done.\n", (unsigned long) pthread_self(), id);
    free(arg);  // free if dynamically allocated
}

int main() {
    thread_pool* pool = create_pool(4);
    if (pool == NULL) {
        printf("Failed to create pool!\n");
        return 1;
    }

    for (int i = 0; i < 10; i++) {
        int* arg = malloc(sizeof(int));
        *arg = i;
        push_task(pool, print_task, arg);
    }

    //  sleep(5);
    destroy_pool(pool, MAX_THREAD_COUNT);
    printf("All tasks done, pool destroyed.\n");
    return 0;
}
