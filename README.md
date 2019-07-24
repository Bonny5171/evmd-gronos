# Gronos

Serviço _cron_ para agendar _worker jobs_ no Faktory Server

### Requisitos

-   Go v1.12.0+ ([download e instruções de instalação e configuração](https://golang.org/dl))

### Criar diretório dos projetos da Everymind no workspace do Go

```bash
$ mkdir -p $GOPATH/src/bitbucket.org/everymind
```

### Clonar este repositório

```bash
$ cd $GOPATH/src/bitbucket.org/everymind
$ git clone https://seu-usuario-git@bitbucket.org/everymind/gronos.git
```

### Clonar o repositório gopkgs

```bash
$ cd $GOPATH/src/bitbucket.org/everymind
$ git clone https://seu-usuario-git@bitbucket.org/everymind/gopkgs.git
```

### Baixar dependências de repositórios externos

```bash
$ cd $GOPATH/src/bitbucket.org/everymind
$ make get
```
