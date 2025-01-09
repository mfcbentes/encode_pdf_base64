# Encode PDF to Base64

Este projeto é uma aplicação Go que converte arquivos PDF em strings Base64.

## Pré-requisitos

- Docker
- Docker Compose

## Configuração

1. Clone o repositório:

   ```sh
   git clone https://github.com/mfcbentes/encode_pdf_base64.git
   cd encode_pdf_base64
   ```

2. Crie um arquivo [.env](http://_vscodecontentref_/1) na raiz do projeto com base no arquivo [EXEMPLE.env](http://_vscodecontentref_/2) e adicione suas variáveis de ambiente:

   ```env
   DB_USER=seu_usuario
   DB_PASSWORD=sua_senha
   DB_PORT=sua_porta
   DB_HOST=seu_host
   DB_SERVICE=seu_servico
   ```

3. Coloque seus arquivos PDF no diretório [input](http://_vscodecontentref_/3).

## Construção e Execução

1. Construa e inicie os contêineres Docker:

   ```sh
   docker-compose up --build
   ```

2. A aplicação estará disponível em `http://localhost:8081`.

## Estrutura do Projeto

- [main.go](http://_vscodecontentref_/4): Código fonte principal da aplicação.
- [Dockerfile](http://_vscodecontentref_/5): Dockerfile para construir a imagem da aplicação.
- [docker-compose.yaml](http://_vscodecontentref_/6): Arquivo de configuração do Docker Compose.
- [input](http://_vscodecontentref_/7): Diretório para armazenar os arquivos PDF de entrada.
- [.env](http://_vscodecontentref_/8): Arquivo de variáveis de ambiente (não incluído no repositório).
- [EXEMPLE.env](http://_vscodecontentref_/9): Exemplo de arquivo de variáveis de ambiente.

## Uso

Envie uma requisição para o endpoint da aplicação para converter um arquivo PDF em uma string Base64.

## Licença

Este projeto está licenciado sob a Licença MIT. Veja o arquivo LICENSE para mais detalhes.
