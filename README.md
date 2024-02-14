## How works:

- Initialy are declare two globals variables ctx and client, after created a type RateLimiterConfig that is a struct.

- func init(), is used for load variables .env and client of the redis.

- In the func main init gin after created local variable ipConfig passing the configs of the ip pattern and below we use rateLimiterMiddleware passing a local variable ipConfig. The same happen with tokenConfig below and finally I used r.GET for set up the route and r.Run for set up the port ":8080".

- The func handleRequest is used only take a message when status is ok.

- The func rateLimiterMiddleware make role of intermediate the quantity of the requests. In your param pass the configs ip and token set ups. It logic implemented is when config.UseIP is true the key is set up the same below for token. After is used a func getLimit passing the key and limit and wait a limit or error. Then we have a count that is used for the quantity of requests and client.Expire is configured the block time, when the number of request pass for the maximum configured throw the message "You have reached the maximum number of requests or actions allowed within a certain time frame".

- It func getLimit serves for get the limits set up in the .env or in database of redis. First get limit in redis, case don't have data this is insert in client.Set and return defaultLimit after compare if limitStr is different of empty string if have a value it convert of the string for int using strconv.Atoi(limitStr) and save data in limit. Other comparison if limit is equal defaultLimit return limit, is because for in case of the be different the new value passed in defaultLimit is save in redis in this case and return the value defaultLimit

- In the func getEnvInt is used for get values in file .env

## How to test:

- set up in file .env the variables IP_LIMIT, TOKEN_LIMIT to configure maximum values for requests. When the values aren't pass, by default IP_LIMIT=10 and TOKEN_LIMIT=100.

- Run the commands: `go run main.go` and `docker compose up -d redis`

- Open in navigator the address: `http://localhost:8080/api/resource`
  - When you use for navigator the the rate limit is used for ip
- For test with token open the client request as Postman and set up in header request API_TOKEN. We put a public workspace for tests in Postman, but you to use ever other.
  For Postman:

  [<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://god.gw.postman.com/run-collection/11060415-1638a203-8d64-43c1-a54d-3a952f52acf7?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D11060415-1638a203-8d64-43c1-a54d-3a952f52acf7%26entityType%3Dcollection%26workspaceId%3D7696cb39-b791-4810-a314-093dfe2d4ca0)

## Automatic tests

- Run `go test`

<!-- Objetivo: Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

_Descrição_ : O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
API_KEY: <TOKEN>
As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

## Requisitos:

O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
O sistema deve responder adequadamente quando o limite é excedido:
Código HTTP: 429
Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
A lógica do limiter deve estar separada do middleware.

## Exemplos:

Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.
Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

## Dicas:

Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.

## Entrega:

O código-fonte completo da implementação.
Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
O servidor web deve responder na porta 8080. -->
