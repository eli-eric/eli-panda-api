MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: 'a8c0e76d-a666-4a82-9fcb-d07b99734b81', code: '0', name: 'L1 Front end (OSCILLATOR and OPA1-3)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '87341236-e829-45de-a0c5-77ce2b336d76', code: '1', name: 'L1 KUBIK 1030 nm amplifier including compressor and SHG'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '99d9225d-7429-4ca1-9431-5669f2e2f33c', code: '2', name: 'L1 OPA4  stage'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '7ffee816-081c-489d-89fe-a2f6ff605d47', code: '3', name: 'L1 ALLEGRA IRT (image relay from OPA4 to OPA5) and FSYNC IRT'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '6deb7286-480a-4f65-941d-bf92e80c5d1d', code: '4', name: 'L1 KUBA 1030 nm amplifier including compressor and SHG'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: 'b075b20e-7bef-4eff-afdc-a86932aa3f8e', code: '5', name: 'L1 Pump laser PL3 (DIRA1 +  B3 compressor incl. beam diagnostics)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: 'a33b245d-399b-4068-b4f9-370f950e4a00', code: '6', name: 'L1 Pump laser PL4 (DIRA2 +  B1 compressor incl. beam diagnostics)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: 'fc013dbe-e89e-435c-9e54-ce4ef78b25be', code: '7', name: 'L1 Pump laser PL5 (DIRA3 + multipass + B2 compressor incl. beam diagnostics)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '1ad4a0f3-7b1f-4251-9c59-ded4a330ea0f', code: '8', name: 'L1 OPCPA (OPA5-7 + diagnostics tower)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '16444493-9e2d-4a9b-b84f-e8a933fc2595', code: '9', name: 'L1 CMC (chirped mirror compressor)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '2cc9eed1-f3c3-4705-bd07-00992ca1ddda', code: 'A', name: 'L1 injector'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '7745758f-7fef-4272-a569-0ab4563e8812', code: 'F1', name: 'FSYNC EMILKA pump laser (regen amplifier and compressor)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: 'abe67d0a-8963-4fa8-9b8e-eb19397e1587', code: 'F2', name: 'FSYNC supercontinuum '}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '98579d9f-a44f-49dc-b2c3-4f5aca738a28', code: 'F3 ', name: 'FSYNC OPA'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '293dd427-169a-425e-80d9-cc42be2c4cf3', code: 'F4 ', name: 'FSYNC JITKAS (optical cross correlators)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '36b44acd-4265-48bf-accc-2fee07b8c2cf', code: 'F5', name: 'FSYNC beam combining setup'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: 'cb22cf42-6bf5-4fe9-88e2-c567f8600d71', code: 'F6', name: 'FSYNC image realy transport from OPA to CMC compressor)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '09bc76b8-7bae-441f-913c-0fa470a29e9b', code: '01', name: 'L1 laser system (L1 Allegra and FSYNC)'}) CREATE (sz:Zone{uid: '7dc00423-4305-4def-9a47-1bf96abb5b07', code: 'F7', name: 'FSYNC CMC compressor'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '388d2b6a-f920-4f9a-82ed-fef3324c6851', code: '02', name: 'L2 DUHA laser'}) CREATE (sz:Zone{uid: '6d1054f8-e068-43f0-a8a6-cf7e6bf6eb04', code: '0', name: 'DUHA BBFE pump laser(including oscillator, and ZBYNA with compressor )'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '388d2b6a-f920-4f9a-82ed-fef3324c6851', code: '02', name: 'L2 DUHA laser'}) CREATE (sz:Zone{uid: '6300fd00-bb80-44a0-9569-17b88b886b06', code: '1', name: 'DUHA BBFE near IR'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '388d2b6a-f920-4f9a-82ed-fef3324c6851', code: '02', name: 'L2 DUHA laser'}) CREATE (sz:Zone{uid: '6688f6a4-b5de-4094-b37d-6c9a4c78eccd', code: '2', name: 'DUHA BBFE mid IR'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '388d2b6a-f920-4f9a-82ed-fef3324c6851', code: '02', name: 'L2 DUHA laser'}) CREATE (sz:Zone{uid: 'cabd67fe-28ee-4f63-a14c-8ee9920f5247', code: '3', name: 'DUHA ns Stretcher'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '388d2b6a-f920-4f9a-82ed-fef3324c6851', code: '02', name: 'L2 DUHA laser'}) CREATE (sz:Zone{uid: '18f7419a-ab62-4555-936c-1bcd240d2fb1', code: '4 ', name: 'DUHA High energy OPCPA (near IR)'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '388d2b6a-f920-4f9a-82ed-fef3324c6851', code: '02', name: 'L2 DUHA laser'}) CREATE (sz:Zone{uid: 'a0b8de4a-57c6-454b-b0c6-a93c88686b21', code: '5', name: 'DUHA 10 J from end'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '388d2b6a-f920-4f9a-82ed-fef3324c6851', code: '02', name: 'L2 DUHA laser'}) CREATE (sz:Zone{uid: '71642441-730e-4667-9c85-a82f36d23e0f', code: '6', name: 'DUHA 10 J amplifier + SHG + transport'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' })
MERGE (z:Zone{uid: '388d2b6a-f920-4f9a-82ed-fef3324c6851', code: '02', name: 'L2 DUHA laser'}) CREATE (sz:Zone{uid: 'f7ed7dc3-cac6-4963-8f0a-5e1b80b9bf81', code: '7', name: 'DUHA Compressor and injector'}) CREATE(z)-[:HAS_SUBZONE]->(sz) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: '2f45c3c3-1833-4bd6-b073-0ac94aab56ec', code: '03', name: 'L3 laser'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: 'c0873468-d49f-45d8-b9ad-beb83d9c9772', code: '04', name: 'L4 laser'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: 'ff46e261-4390-491d-b784-119284acdf3a', code: '05', name: 'L3 Beam transport'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: '04945866-a2c1-4a56-9b97-0cb22c659cad', code: '06', name: 'E1 experiments'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: '620af539-4d0f-4ce9-a18f-bef4f9e11a38', code: '07', name: 'E2 experiments '}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: 'a2d5c95f-f415-4947-b962-ad742216278d', code: '08', name: 'E3 experiments'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: '258fb941-1c02-4c1c-ba8d-50cf0deb85d4', code: '09', name: 'E4 Experiments'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: 'ae935a2a-512c-40a4-801c-bf72d3d7028d', code: '10', name: 'E5 experimental hall'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: '490ce3cc-4e1d-4eef-9e47-60ae3eb7398b', code: '11', name: 'E6 experimental hall'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: '685eb6b9-b238-471a-a617-2412c2a83c65', code: '12 ', name: 'Facility systems'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: '53690c65-17cd-4799-be2c-ebc63bda421f', code: '13', name: 'E5 LUIS experiment'}) CREATE(f)-[:HAS_ZONE]->(z);

MATCH (f:Facility{ code: 'B' }) CREATE(z:Zone{uid: 'afecb84f-9fb4-4baa-805d-2c3093cc6528', code: '14', name: 'E5 ELBA experiment'}) CREATE(f)-[:HAS_ZONE]->(z);
