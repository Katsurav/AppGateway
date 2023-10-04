# TTECAppGateway
 TTEC KIA Application Gateway

The Application Gateway listens on Port 8080

Supported Actions

    GET /health
    GET /config
    POST /query

Example Query:

POST a JSON message like below:

{ "user": "smlsr",
  "query": "how to fix audio?",
  "debug": false,
  "pre": "",
  "post": "",
  "engines": [ "gpt-35-turbo" ] }

    user:     Optional for now
    query:    Mandatory -> The question you would like answered
    debug:    Optional for now
    pre:      If you would like to PREPEND all LLM queries with a Prompt modification enter it here, using "!" as first character will DELETE the existing PREPEND prompt defined in the AppGateway
    post:     If you would like to APPEND all LLM queries with a Prompt modification enter it here, using "!" as first character will DELETE the existing APPEND prompt defined in the AppGateway
    engines:  gpt-35-turbo, text-davinci-003, palm are options (can be multiple)


