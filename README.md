# rest api for finding openapi json/yaml spec links using the supplied query term

Web API for finding openapi json/yaml spec links using the supplied query term. The query term is searched 
in bing web search for open api specification links. The resultant urls are crawled using colly package for 
json/yaml links

## main.go

usage : http://127.0.0.1:8000/?qt=google pay

loads the config.json and initialises a Web server 
handler accepts a query term for searching
call the search function
call the crawl function with the list of urls from search
output the links to yaml/json files found in json format ordered by rank and depth

output json structure
[{"LinkTitle":"OpenAPI Payment Processor Service - JSON","Referrer":"https://developers.google.com/standard-payments/v1/payment-update-service-api/open-api-spec","URL":"https://developers.google.com/standard-payments/v1/payment-processor-service-api/open-api-spec.json","Depth":1,"Rank":0,"Baseurl":"https://developers.google.com/standard-payments/v1/payment-update-service-api/open-api-spec"},{"LinkTitle":"OpenAPI Payment Processor Service - JSON","Referrer":"https://developers.google.com/standard-payments/v1/payment-processor-service-api/open-api-spec","URL":"https://developers.google.com/standard-payments/v1/payment-processor-service-api/open-api-spec.json","Depth":1,"Rank":1,"Baseurl":"https://developers.google.com/standard-payments/v1/payment-processor-service-api/open-api-spec"}]

## websearch.go

Search using bing web search api with the given query term 
Output the list of urls found in an array

## collycrawl.go

crawl the urllist and output the list of json/yaml file links in a json format

## config.json

configs for 
* web server host and port
* bing web search api endpoint , token, resultlimit and search suffix 
* colly maxdepth, parallelism, excluded url patterns, excluded file extensions
