# Rate Limiting com Redis em Go

Este projeto implementa uma solução de **rate limiting** para uma aplicação GO utilizando o Redis como backend para armazenar o estado dos limites de requisições.

## Descrição

O objetivo do rate limiting é restringir o número de requisições que um usuário pode fazer em um determinado período de tempo. Neste projeto, foi implementado um mecanismo de rate limiting com o Redis, onde a aplicação permite um número máximo de requisições por um usuário, e após atingir esse limite, as requisições subsequentes são bloqueadas até o tempo limite expirar.

## Tecnologias Utilizadas

- **Go (Golang)**: Linguagem de programação para implementar a API.
- **Redis**: Banco de dados em memória utilizado para armazenar e controlar o número de requisições feitas.
- **go-redis**: Biblioteca Go para interação com o Redis.

## Funcionalidade

### Fluxo do Rate Limiting:

1. Quando uma requisição é feita à API, o sistema verifica no Redis o número de requisições feitas por um usuário específico no período de tempo configurado.
2. Se o limite de requisições não for alcançado, a requisição é processada normalmente.
3. Caso o limite de requisições tenha sido atingido, a requisição é rejeitada com uma resposta de erro informando que o limite foi excedido.
4. A contagem de requisições é resetada após o tempo limite configurado.
5. A aplicação por parão irá fazer o rate limit pelo IP do usuário através das variáveis de ambiente `LIMIT_PER_SECOND` e  `BLOCK_DURATION_PER_SECOND` que por padrão tem os valores de 5 e 10 respectivamente (podendo ser configurado nas variaveis de ambiente), porém, caso for enviado no header da requisição um token `API_KEY` que pode onde corresponda com o valor da variável de ambiente `CUSTOM_TOKEN_REQUESTS` será considerado a quantidade de requests configurada no `CUSTOM_TOKEN_REQUESTS_VALUE`

### Como Funciona:

1. **Chave de Rate Limit**: Para cada usuário, a chave no Redis é composta pelo identificador do usuário (exemplo: endereço IP ou ID de autenticação) e o período de tempo (exemplo: 1 minuto, 1 hora).
   
2. **Definição de Limite**: O número máximo de requisições permitidas para o usuário no período especificado é armazenado e consultado.
   
3. **Expiração de Chave**: Cada chave criada para controlar as requisições tem um tempo de expiração, que corresponde ao período durante o qual o usuário pode fazer requisições.

## **Instruções de Execução**
### **Rodando com Docker Compose**
1. Clone o repositório:
   ```sh
   git clone https://github.com/brunoofgod/goexpert-lesson-7.git
   cd goexpert-lesson-7
   ```

2. Inicie o projeto com Docker Compose:
   ```sh
   docker-compose up -d
   ```

