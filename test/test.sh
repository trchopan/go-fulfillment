# curl -X POST https://go-fulfillment-ywozb2wfqq-uc.a.run.app/createFulfilment\
  # -d '{"title":"what a dash", "body":"Nuuuu batachi"}'

# curl -X PUT http://localhost:8000/updatePost/1000000 -d '{"title":"naniii", "body":"wooooootaasdfad asdf asfs dfasdfa f"}'

curl -X POST http://localhost:8080/fulfillment \
  -d '{"title":"what a dash", "body":"Nuuuu batachi"}'
