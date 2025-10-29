#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<string.h>
#include<stdbool.h>

const int MAX_KEY = 15;
const int MAX_VALUE = 100;
const int MAX_METHOD = 10;
const int MAX_PATH = 100;
const int MAX_PROTOCOL = 20;
const int MAX_HEADER = 10;
const int MAX_BUF = 1024;

typedef enum {
	METHOD = 0,
	HEADER,
	BODY
} parse_state;

typedef struct header {
	char key[MAX_KEY];
	char value[MAX_VALUE];
} header;

typedef struct {
	char method[MAX_METHOD];
	char path[MAX_PATH];
	char protocol[MAX_PROTOCOL];
	header hdr[MAX_HEADER];
	int hdrCount;
	char *body;
	int err;
} http_request;

void parse_line(http_request* r, char* line, parse_state state) {
	switch (state) {
	case METHOD:
		sscanf(line, "%s %s %s", r->method, r->path, r->protocol);
		break;
	case HEADER: {
		int idx = r->hdrCount;
		int breakPoint = strcspn(line, ":");	// index of ':'

		// extract key
		strncpy(r->hdr[idx].key, line, sizeof(char) * breakPoint);
		r->hdr[idx].key[breakPoint] = '\0';

		// remove space
		while (line[++breakPoint] == ' ') {}

		// extract value
		strncpy(r->hdr[idx].value, line + breakPoint, sizeof(char) * (strlen(line) - breakPoint));
		r->hdr[idx].value[ (strlen(line) - breakPoint)] = '\0';

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
	for (int i = 0, j = 0; i <= strlen(raw_request); i++) {
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
		if (r->body) free(r->body);
		free(r);
	}
}

int main() {

	char raw_request1[] =
	    "GET /index.html?user=labib&lang=en HTTP/1.1\r\n"
	    "Host: localhost:8080\r\n"
	    "User-Agent: curl/8.0\r\n"
	    "Accept: text/html\r\n"
	    "Connection: keep-alive\r\n"
	    "\r\n";

	char raw_request2[] =
	    "POST /submit-form HTTP/1.1\r\n"
	    "Host: localhost:8080\r\n"
	    "User-Agent: curl/8.0\r\n"
	    "Content-Type: application/x-www-form-urlencoded\r\n"
	    "Content-Length: 27\r\n"
	    "\r\n"
	    "username=labib&lang=en&age=20";

	char raw_request3[] =
	    "GET / HTTP/1.1\r\n"
	    "Host: 127.0.0.1:8080\r\n"
	    "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:144.0) Gecko/20100101 Firefox/144.0\r\n"
	    "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n"
	    "Accept-Language: en-US,en;q=0.5\r\n"
	    "Accept-Encoding: gzip, deflate, br, zstd\r\n"
	    "Connection: keep-alive\r\n"
	    "Upgrade-Insecure-Requests: 1\r\n"
	    "Sec-Fetch-Dest: document\r\n"
	    "Sec-Fetch-Mode: navigate\r\n"
	    "Sec-Fetch-Site: none\r\n"
	    "Sec-Fetch-User: ?1\r\n"
	    "Priority: u=0, i\r\n"
	    "\r\n";

	PRINT(raw_request1);
	PRINT(raw_request2);
	PRINT(raw_request3); // If you see any error, just change the max values ::)

	return 0;
}
