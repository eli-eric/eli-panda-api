// Remove PublishFormat nodes
MATCH (pf:PublishFormat) DETACH DELETE pf;

// Remove ConferenceScope nodes
MATCH (cs:ConferenceScope) DETACH DELETE cs;
