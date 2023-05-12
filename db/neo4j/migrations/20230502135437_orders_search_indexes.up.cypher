CREATE FULLTEXT INDEX searchIndexOrders IF NOT EXISTS FOR (n:Order) ON EACH [n.name, n.orderNumber, n.requestNumber, n.contractNumber, n.notes]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'standard'
  }
};

CREATE INDEX ordersUidIndex IF NOT EXISTS FOR (o:Order) ON o.uid;