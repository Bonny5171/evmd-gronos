# Gronos

Serviço _cron_ foi criado para agendar _worker jobs_ no Faktory Server

### Requisitos

-   Go v1.13.0+ ([download e instruções de instalação e configuração](https://golang.org/dl))

### Criar diretório dos projetos da Everymind no workspace do Go

```bash
$ mkdir -p $GOPATH/src/bitbucket.org/everymind
```

### Clonar este repositório

```bash
$ cd $GOPATH/src/bitbucket.org/everymind
$ git clone https://seu-usuario-git@bitbucket.org/everymind/evmd-gronos.git
```

### Clonar o repositório evmd-golib

```bash
$ cd $GOPATH/src/bitbucket.org/everymind
$ git clone https://seu-usuario-git@bitbucket.org/everymind/evmd-golib.git
```

### Clonar o repositório gforce

```bash
$ cd $GOPATH/src/bitbucket.org/everymind
$ git clone https://seu-usuario-git@bitbucket.org/everymind/gforce.git
```

### Baixar dependências de repositórios externos

```bash
$ cd $GOPATH/src/bitbucket.org/everymind
$ go mod download
```
