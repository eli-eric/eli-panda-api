// -------- Nodes (labels + properties) --------
CALL {
  CALL db.schema.nodeTypeProperties()
  YIELD nodeLabels, propertyName, propertyTypes
  UNWIND nodeLabels AS label
  WITH label, propertyName, propertyTypes
  WITH label,
       collect(DISTINCT {
         name: propertyName,
         type: toString(coalesce(head(propertyTypes), 'ANY'))
       }) AS props
  RETURN collect(DISTINCT {label: label, properties: props}) AS nodes
}

// -------- Relationship endpoints (type + start + end) from data --------
CALL {
  // List all relationship types first
  CALL db.relationshipTypes() YIELD relationshipType
  WITH relationshipType AS type
  // For each type, collect distinct start/end label sets from existing relationships
  CALL {
    WITH type
    OPTIONAL MATCH (a)-[r]->(b)
    WHERE type(r) = type
    WITH collect(DISTINCT labels(a)) AS ssets,
         collect(DISTINCT labels(b)) AS esets
    RETURN apoc.coll.toSet(apoc.coll.flatten(ssets)) AS start,
           apoc.coll.toSet(apoc.coll.flatten(esets)) AS end
  }
  RETURN collect({type: type, start: start, end: end}) AS rel_ends
}

// -------- Relationship properties (per type) --------
CALL {
  CALL db.schema.relTypeProperties()
  YIELD relType, propertyName, propertyTypes
  WITH relType,
       collect(DISTINCT {
         name: propertyName,
         type: toString(coalesce(head(propertyTypes), 'ANY'))
       }) AS props
  RETURN collect({type: relType, properties: props}) AS rel_props
}

// -------- Merge rel endpoints + props via lookup --------
WITH nodes, rel_ends, rel_props,
     apoc.map.fromPairs([rp IN rel_props | [rp.type, rp.properties]]) AS relPropsByType
WITH nodes,
     [rel IN rel_ends |
       {
         type: rel.type,
         start: rel.start,        // e.g., ["Person"]
         end: rel.end,            // e.g., ["Movie"]
         properties: coalesce(relPropsByType[rel.type], [])
       }
     ] AS relationships
RETURN apoc.convert.toJson({nodes: nodes, relationships: relationships}) AS schema_json;
