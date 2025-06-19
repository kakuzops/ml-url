   ### Melhorias
   - Implementar rate limiting
   - Adicionar autenticação JWT
   - Adicionar validação de URLs maliciosas


   nivel de escala da aplicação ->
 quantos caracteres o shortner aguenta-> entre 6 e 8 caracteres
   |> 68.719.476.736
   
   oque acontece quando tiver 100m de request o prometheus consegue resolver isso?
   
   |> Reduzir a cardinalidade das métricas
   |> Implementar rate limiting no endpoint de métricas (com 100m de requests  isso pode gerar milhoes de info diferentes) custo e problema a longo prazo

   quais as urls mais chamadas
   (post e get)

   qual seria o deployment disso?
      |> seria um deployment com kubernets usando load balancer para o nivel meli(50m de request min)
      |> NGINX para gerenciamento de trafego externo.
      |> HPA dos pods para escalar horizontalmente
      |> com replica de banco de dados para escrita e leitura(1 primario e 2 replicas)
      |> um pooling com pgbouncer
      |> cache no redis com 2-3 nos
   


   - diagrama da arquitetura


capacidade de analise
colocar que eu assumi tal produto ia levar x tempo ou x requests








