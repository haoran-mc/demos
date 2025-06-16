docker run -d --name=ran_neo4j -p 7474:7474 -p 7687:7687 -v /Users/haoran/neo4j/import:/var/lib/neo4j/import neo4j
mkdir -p ~/neo4j/import
cp ./*.csv ~/neo4j/import



加载诗人 (poet.csv)
LOAD CSV WITH HEADERS FROM 'file:///poet.csv' AS row
CREATE (:Poet {poet_id: toInteger(row.`:ID(poet-id)`), name: row.name});



加载诗 (poem.csv)
LOAD CSV WITH HEADERS FROM 'file:///poem.csv' AS row
CREATE (:Poem {poem_id: toInteger(row.`:ID(poem-id)`), name: row.name});



加载诗人与诗的关系 (edge.csv)
LOAD CSV WITH HEADERS FROM 'file:///edge.csv' AS row
MATCH (p:Poet {poet_id: toInteger(row.`:START_ID(poet-id)`)})
MATCH (m:Poem {poem_id: toInteger(row.`:END_ID(poem-id)`)})
CREATE (p)-[:CREATED {relation: row.relation}]->(m);





查看所有诗人 (Poet)
MATCH (p:Poet) RETURN p LIMIT 10;



查看所有诗 (Poem)
MATCH (m:Poem) RETURN m LIMIT 10;



查看诗人与诗的关系
MATCH (p:Poet)-[r:CREATED]->(m:Poem) RETURN p, r, m LIMIT 10;





查询某个诗人的作品
例如，查看 李白 的作品：
MATCH (p:Poet {name: '李白'})-[:CREATED]->(m:Poem)
RETURN p.name AS Poet, m.name AS Poem;



查询某首诗的作者
例如，查询《将进酒》的作者：
MATCH (p:Poet)-[:CREATED]->(m:Poem {name: '将进酒'})
RETURN m.name AS Poem, p.name AS Poet;
