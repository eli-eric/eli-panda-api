CREATE FULLTEXT INDEX searchIndexSystems IF NOT EXISTS FOR (n:System) ON EACH [n.name, n.description, n.systemCode, n.systemAlias]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'standard'
  }
};
