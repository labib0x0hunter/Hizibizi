#include<stdio.h>
#include<stdlib.h>
#include<signal.h>
#include<unistd.h>
#include<string.h>
#include<stdbool.h>
#include<pthread.h>
#include<sys/socket.h>
#include<netinet/in.h>

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

const int MAX_KEY = 50;
const int MAX_VALUE = 1024;
const int MAX_METHOD = 10;
const int MAX_PATH = 512;
const int MAX_PROTOCOL = 20;
const int MAX_HEADER = 50;
const int MAX_BUF = 1024;

typedef enum {
    METHOD = 0,
    HEADER,
    BODY
} parse_state;

typedef struct header {
	char key[MAX_KEY];
	char value[MAX_VALUE];
    int key_len;
    int value_len;
} header;

typedef struct {
	char method[MAX_METHOD];
	char path[MAX_PATH];
	char protocol[MAX_PROTOCOL];
	header hdr[MAX_HEADER];
	int hdrCount;
	char *body;
	int err;
    int method_len;
    int path_len;
    int protocol_len;
} http_request;

void parse_line(http_request* r, char* line, parse_state state) {
	switch (state) {
		case METHOD:
            {
                int i = 0;
                int cnt = 0;
                for(int j = 0; j < (int) strlen(line); j++) {
                    if (line[j] == ' ') { 
                        if (cnt == 0) {
                            r->method[i] = '\0';
                            r->method_len = i;
                        } else if (cnt == 1) {
                            r->path[i] = '\0';
                            r->path_len = i;
                        }
                        cnt++;
                        i = 0;
                    }
                    while (j <= (int) strlen(line) && line[j] == ' ') { j++; }
                    if (cnt == 0) {
                        r->method[i++] = line[j];
                    } else if (cnt == 1) {
                        r->path[i++] = line[j];
                    } else {
                        r->protocol[i++] = line[j];
                    }
                }
                r->protocol[i] = '\0';
                r->protocol_len = i;
                break;
            }
		case HEADER: {
			int idx = r->hdrCount;
			int breakPoint = strcspn(line, ":");	// index of ':'

			// extract key
			strncpy(r->hdr[idx].key, line, sizeof(char) * breakPoint);
			r->hdr[idx].key[breakPoint] = '\0';
            r->hdr[idx].key_len = breakPoint;

			// remove space
			while(line[++breakPoint] == ' ') {}

			// extract value
			strncpy(r->hdr[idx].value, line + breakPoint, sizeof(char) * (strlen(line) - breakPoint));
			r->hdr[idx].value[ (strlen(line) - breakPoint)] = '\0';
            r->hdr[idx].value_len = strlen(line) - breakPoint;

			r->hdrCount++;
			break;
		}
		case BODY:
			r->body = (char*) malloc(sizeof(char) * (strlen(line) + 1)); // what if malloc() fails ???
			if (r->body == NULL) {
				r->err = -1;
				return;
			}
			strncpy(r->body, line, strlen(line));
			r->body[strlen(line)] = '\0';
			break;
	}
}

http_request* parse_http_request(char* raw_request) {
	http_request* req = (http_request*) malloc(sizeof(http_request));
	if (req == NULL) {
		return NULL;
	}
	// set header count = 0
	req->hdrCount = 0;
	req->err = 0;

	// buffer to store line
	char buffer[MAX_BUF];
	parse_state state = METHOD;	// first state 
	for (int i = 0, j = 0; i <= (int) strlen(raw_request); i++) {
		if (req->err != 0) {
			break;
		}
		switch (raw_request[i]) {
		case '\r':
			if (raw_request[i + 1] == '\n') {
				i++;
			}
			buffer[j] = '\0';
			j = 0;
			if (state == METHOD) {
				parse_line(req, buffer, METHOD);
				state = HEADER;
			} else if (state == HEADER) {
				if (strlen(buffer) == 0) {	// copy the remaining of the body and return..
					state = BODY;
					int remaining = strlen(raw_request) - (i + 1);
					if (remaining > 0) {
						req->body = malloc(remaining + 1);
						if (req->body == NULL) {
							req->err = -1;
							break;
						}
						strcpy(req->body, raw_request + i + 1);
					}
					return req;
				} else {
					parse_line(req, buffer, HEADER);
				}
			}
			break;
		case '\0':
			buffer[j] = '\0';
			if (strlen(buffer) != 0) {
				parse_line(req, buffer, BODY);
			}
			break;
		default:
			buffer[j++] = raw_request[i];
			break;
		}
	}
	return req;
}

void PRINT(char* raw_request) {
	http_request* r = parse_http_request(raw_request);
	if (r == NULL) {
		return;
	} else if (r->err == -1) {
		return;
	}
	printf("METHOD = %s\n", r->method);
	printf("PATH = %s\n", r->path);
	printf("PROTOCOL = %s\n", r->protocol);
	for (int i = 0; i < r->hdrCount; i++) {
		printf("KEY = %s , VALUE = %s\n", r->hdr[i].key, r->hdr[i].value);
	}
	printf("BODY = %s\n\n", r->body);
	if (r) {
		if(r->body) free(r->body);
		free(r);
	}
}

typedef struct client {
	int fd;
} client;

// handles client
void handle_client(void* arg) {
	client* c = (client*) (arg);
	char request[1024];
    int n = recv(c->fd, request, sizeof(request), 0);
	request[n] = '\0';
	// printf("%s\n", request);
	PRINT(request);

    http_request* r = parse_http_request(request);

    // char response[] = ;

	n = send(c->fd, r->path, sizeof(char) * r->path_len, 0);

    // printf("DEBUG :: %d  -  %d \n", r->path_len, n);
    if (n != r->path_len) {
        // partical send.........
    }

	close(c->fd);
	free(c);
}

thread_pool* pool;
int server_fd;
volatile sig_atomic_t shutdown_signal;

void shut_down_server(int sig) {
    printf("%d signal... shutdown..", sig);
    shutdown_signal = 1;
    close(server_fd);
}

int main() {

    shutdown_signal = 0;

	pool = create_pool(4);
    if (pool == NULL) {
        printf("Failed to create pool!\n");
        return 1;
    }

	// create a socket endpoint
	server_fd = socket(AF_INET, SOCK_STREAM, 0);
	if (server_fd < 0) {
		perror("socket create");
		return 1;
	}
	printf("%d\n", server_fd);

    // handle ctrl + c 
    signal(SIGINT, shut_down_server);

	struct sockaddr_in addr;
	bzero(&addr, sizeof(addr));	// fill with zero

	addr.sin_family = AF_INET;	// ipv4
	addr.sin_port = htons(8080);	// port -> endianess handles
	addr.sin_addr.s_addr = INADDR_ANY;	// accept any local interface

	if (bind(server_fd, (struct sockaddr*)&addr, sizeof(addr)) < 0 ) { perror("bind"); return 1; };	// bind to port 127.0.0.1:8080
	if (listen(server_fd, 10) < 0) { perror("listen"); return 1; }; // 10 backlog queue, why use backlog queue ? maximum number of pending connections waiting to be accepted

	while (!shutdown_signal) {
		int c_fd = accept(server_fd, NULL, NULL);	// blocks until a new client connects ..
		if (c_fd < 0) {
			continue;
		}
		client* c = (client*) malloc(sizeof(client));
		if (c == NULL) { continue; }
		c->fd = c_fd;
		push_task(pool, handle_client, c);	// push to thread-pool
	}

    if (!shutdown_signal) {
	    close(server_fd);
    }
    destroy_pool(pool, MAX_THREAD_COUNT);

	return 0;
}
