Adicionar latência, exemplo do GPT:

Além do status HTTP, mede latência:

inicio := time.Now()
resp, err := http.Get(site)
duracao := time.Since(inicio)

--------------

Retry automático, se um site falhar tenta novamente 3x e espera 1s entre as tentativas;

--------------

Melhorar a concorrência com goroutines:

monitorar centenas de sites simultaneamente;
limitar workers;
criar pool de workers;

Exemplo de arquitetura do GPT:
Sites -> Channel -> Workers -> Resultados

--------------

Criar o software com um front para cadastrar novos sites para monitoramento:
opção para cadastrar novos sites, alimenta o banco, próxima atualização já considera o site.

--------------

Verificar outras informações como:
verificar o conteúdo do 'GET';
validar JSON do GET;
validar headers;
validar tempo máximo, onde resposta > 2s = ALERTA;

--------------

Criar tabelas de bancos para:
cadastro de sites;
registro de status (ativo/inativo, quando, código http);
tempo de resposta, criar um banco para monitorar esses tempos de respostas;
timestamp.

--------------

Criar alertas reais:
quando o site cair, enviar slack ou e-mail ou telegram ou discord webhook;