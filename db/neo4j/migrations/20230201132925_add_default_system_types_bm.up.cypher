MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'f8359245-ad17-425b-9e4f-3cf20580db7a', name:'Lens', code: 'L', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'9f59a159-2a09-4460-91e4-6d57fc9ab6b4', name:'Mirror', code: 'M', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'ba8d0150-7320-48ed-b5f6-4d81d0ed652a', name:'Filter', code: 'F', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'f546d0ff-b95b-4af9-b889-604c9be23b08', name:'Crystal', code: 'CR', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'e44eee3b-b179-4e34-b25f-b33f6bd76b39', name:'Beamsplitter', code: 'BS', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'18de2fe5-a0fc-47fb-9e99-e81dfb94a12e', name:'waveplate', code: 'WP', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'2c426de3-19b2-40da-ab13-ef79d4d00ea1', name:'Grating', code: 'G', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'b9af69ac-c3a3-42c2-84c2-e843f5ad377e', name:'Polarizer', code: 'P ', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'dbd4be86-0337-4920-bf8d-c8ad04fb4910', name:'Prism', code: 'PR', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'87e4c772-cc55-42b3-992d-2a74b3cc0fd1', name:'Isolator', code: 'IS', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'def98651-34e1-4d51-a39b-e65619904f1e', name:'Window', code: 'W', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'38400d44-5350-413e-90c0-31142e65e457', name:'Off axis parabola', code: 'OAP', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'50854f14-8734-4c02-90eb-bc3ca2ad6dc3', name:'Optical components' })
MERGE (st:SystemType{ uid:'d3596300-4157-4528-884b-63e058416534', name:'Beamdump', code: 'BD', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d370a073-3cfa-41bc-a98c-65a66df430c6', name:'Laser diagnostics' })
MERGE (st:SystemType{ uid:'0c64442e-ee2f-42e3-b368-02f42518fcb6', name:'Autocorrelator', code: 'AC', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d370a073-3cfa-41bc-a98c-65a66df430c6', name:'Laser diagnostics' })
MERGE (st:SystemType{ uid:'646f4f1f-eddd-4c20-863c-f47727262c4b', name:'Photodiode', code: 'PD', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d370a073-3cfa-41bc-a98c-65a66df430c6', name:'Laser diagnostics' })
MERGE (st:SystemType{ uid:'d3388c81-8bfa-4828-8d7b-64d2e9a580b1', name:'Balanced photodiode ', code: 'BPD', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d370a073-3cfa-41bc-a98c-65a66df430c6', name:'Laser diagnostics' })
MERGE (st:SystemType{ uid:'2b5ab14e-c157-475f-9efd-aacc62884e54', name:'Camera', code: 'C', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d370a073-3cfa-41bc-a98c-65a66df430c6', name:'Laser diagnostics' })
MERGE (st:SystemType{ uid:'7208135c-4401-4b5a-b2e0-77285dfc9dc4', name:'Power meter', code: 'PM', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d370a073-3cfa-41bc-a98c-65a66df430c6', name:'Laser diagnostics' })
MERGE (st:SystemType{ uid:'6f63b2c7-f9ef-4dd5-8a99-8495e14682a2', name:'Energy meter', code: 'EM', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'8f8d85c9-4871-4d29-8087-199f4da92178', name:'Lasers, laser components and devices' })
MERGE (st:SystemType{ uid:'cbd1c773-a1e3-4d01-a84e-c2375ba1a84b', name:'Fiber amplifer', code: 'FA', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'8f8d85c9-4871-4d29-8087-199f4da92178', name:'Lasers, laser components and devices' })
MERGE (st:SystemType{ uid:'5529f852-b297-4305-ab2e-562ddc9971e9', name:'Regenerative amplifier', code: 'RA', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'8f8d85c9-4871-4d29-8087-199f4da92178', name:'Lasers, laser components and devices' })
MERGE (st:SystemType{ uid:'98d23683-2d29-4526-a6ac-990a39be2b19', name:'Laser diodes', code: 'LD', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'8f8d85c9-4871-4d29-8087-199f4da92178', name:'Lasers, laser components and devices' })
MERGE (st:SystemType{ uid:'bf096209-2223-4657-889b-d3c1847318d5', name:'Laser Oscillator', code: 'OSC', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'8f8d85c9-4871-4d29-8087-199f4da92178', name:'Lasers, laser components and devices' })
MERGE (st:SystemType{ uid:'da1af02b-ece5-41b9-9665-4d5cae51d9ff', name:'Pockels cell', code: 'PC', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'a643e22d-f0af-4be1-8069-e17ba4b0dac1', name:'Cooling circuit coponents' })
MERGE (st:SystemType{ uid:'1ba63d9a-5d24-41af-8e42-16855a4a7286', name:'Chiller', code: 'CHL', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'cb3fddac-7982-4e63-b2ca-d63d089a6579', name:'Turbomolecular pump', code: 'TMP', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'adc44705-cfac-4837-b87f-f5592bfd0afd', name:'Vacuum primary pump', code: 'VPP', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'1afe3efa-6a0a-49d7-8487-326431fe5606', name:'Gate valve', code: 'GV', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'aca9c8d9-e623-431e-85c5-06591d6f19a4', name:'Safety gate valve', code: 'SGV', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'7b9c185b-f75b-4828-aad7-ad9dc884a4af', name:'Endstation gate valve with window', code: 'EGV', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'8783fa8c-cfd3-4519-830f-5bfd05166ad9', name:'Pirani vacuum gauge', code: 'PG', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'aa781eef-5d38-4a66-a33c-10718a7cc402', name:'Accucate atmospheric pigani gauge', code: 'APG', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'02b2b92c-4ab3-4c24-948f-838d1e93fe74', name:'Vacuum wide range gauge', code: 'WRG', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'9585a948-4006-4c50-8fa4-22451183cbc4', name:'Venting valve', code: 'VV', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'9711f861-4c46-4009-b9e8-c67a011a3b24', name:'Vacuum roughing valve', code: 'RV', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'d9107f91-bb98-4445-a43b-3b741df4d2f6', name:'Active vacuum devices' })
MERGE (st:SystemType{ uid:'5bad59cb-a52b-4887-b5d7-c4c68935b51d', name:'Vacuum backing valve', code: 'BV', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'4a95531c-d087-479d-82c0-b4da2f3b8ace', name:'Motion devices' })
MERGE (st:SystemType{ uid:'d0104a0c-50fb-498e-9b03-840b3d2be0c0', name:'Picomotor driver', code: 'PMD', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'4a95531c-d087-479d-82c0-b4da2f3b8ace', name:'Motion devices' })
MERGE (st:SystemType{ uid:'fa072c42-55ef-4cec-83f9-06cf1ec8ba4e', name:'Piezo motor driver', code: 'PZMD', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'307033d7-e4a6-4139-aeec-8955455986a8', name:'Controls and drivers' })
MERGE (st:SystemType{ uid:'d2d02f3d-3b00-4c80-ae8a-42c65453bb82', name:'Stepper motor driver', code: 'SMD', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility{ code: 'B' })
MERGE (sg:SystemTypeGroup{ uid:'2850a16e-5158-485c-84f8-e6af44186790', name:'IT devices' })
MERGE (st:SystemType{ uid:'4baea8cd-567f-4607-97b3-f8d4775313a8', name:'Network switch', code: 'NSW', mask: '' })
MERGE (sg)-[:CONTAINS_SYSTEM_TYPE]->(st)
MERGE (sg)-[:BELONGS_TO_FACILITY]->(f);
