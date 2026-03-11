// PublishFormat seed data
MERGE (pf:PublishFormat {code: "P"}) ON CREATE SET pf.uid = apoc.create.uuid(), pf.name = "Print";
MERGE (pf:PublishFormat {code: "E"}) ON CREATE SET pf.uid = apoc.create.uuid(), pf.name = "Online";
MERGE (pf:PublishFormat {code: "C"}) ON CREATE SET pf.uid = apoc.create.uuid(), pf.name = "Physical media (CD, DVD, flash drive)";

// ConferenceScope seed data
MERGE (cs:ConferenceScope {code: "CST"}) ON CREATE SET cs.uid = apoc.create.uuid(), cs.name = "National";
MERGE (cs:ConferenceScope {code: "EUR"}) ON CREATE SET cs.uid = apoc.create.uuid(), cs.name = "European";
MERGE (cs:ConferenceScope {code: "WRD"}) ON CREATE SET cs.uid = apoc.create.uuid(), cs.name = "Worldwide";
