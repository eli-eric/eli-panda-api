MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Atrium', facility:'B' })
MERGE (l2:Location{ name: 'main entrance hall', code: 'A.00.01', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Atrium', facility:'B' })
MERGE (l2:Location{ name: 'stair A1', code: 'A.00.02', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Atrium', facility:'B' })
MERGE (l2:Location{ name: 'bridge', code: 'A.1.01', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Atrium', facility:'B' })
MERGE (l2:Location{ name: 'bridge', code: 'A.1.02', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Atrium', facility:'B' })
MERGE (l2:Location{ name: 'bridge', code: 'A.2.01', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Atrium', facility:'B' })
MERGE (l2:Location{ name: 'stair A1', code: 'A.2.02', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Atrium', facility:'B' })
MERGE (l2:Location{ name: 'personal lift', code: 'A.LI.01', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Cooling Plantroom building', facility:'B' })
MERGE (l2:Location{ name: 'cooling plantroom', code: 'CO.00.01', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Cooling Plantroom building', facility:'B' })
MERGE (l2:Location{ name: 'transformer room', code: 'CO.00.02', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Cooling Plantroom building', facility:'B' })
MERGE (l2:Location{ name: 'power switch room', code: 'CO.00.03', facility:'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office ELI2', code: 'II.02.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office ELI2', code: 'II.02.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office ELI2', code: 'II.02.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office ELI2', code: 'II.02.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office ELI2', code: 'II.02.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office ELI2', code: 'II.02.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Power distribution, measurement and  control room ', code: 'II.02.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Server Room S3', code: 'II.02.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Gas heating room', code: 'II.02.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Kitchen', code: 'II.02.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Server Room S10', code: 'II.02.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Lobby', code: 'II.01.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Reception and staircase', code: 'II.01.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Corridor', code: 'II.01.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Electricasl workshop', code: 'II.01.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Machanical workshop (materail and preparation)', code: 'II.01.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Mechanical workshop (CNCs)', code: 'II.01.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Warehouse Main Area', code: 'II.01.07.1', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'ELI2 building', facility: 'B' })
MERGE (l2:Location{ name: 'Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Warehouse Back Area', code: 'II.01.07.2', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room', code: 'L.00.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L7', code: 'L.00.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room', code: 'L.00.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'North corridor (for personnel access)', code: 'L.00.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L1', code: 'L.00.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'L.00.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L4b laser hall', code: 'L.00.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L3 laser hall', code: 'L.00.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L2 laser hall', code: 'L.00.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L1 Laser hall', code: 'L.00.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'South corridor (for equipment)', code: 'L.00.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L2', code: 'L.00.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L3', code: 'L.00.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L5', code: 'L.00.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L4', code: 'L.00.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L6', code: 'L.00.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'L.00.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'L.00.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L1 oscillator room', code: 'L.00.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L2 control room', code: 'L.00.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L3 control room', code: 'L.00.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L4 control room', code: 'L.00.24', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L1 rack area', code: 'L.00.25', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '00-Laser floor', facility: 'B' })
MERGE (l3:Location{ name: 'L1 control room', code: 'L.00.26', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'platform for vacuum pumps', code: 'L.01.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plenum', code: 'L.01.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plenum', code: 'L.01.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plenum', code: 'L.01.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'platform for vacuum pumps', code: 'L.01.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L1', code: 'L.01.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'platform for vacuum pumps', code: 'L.01.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plenum', code: 'L.01.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L6', code: 'L.01.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plenum', code: 'L.01.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L2', code: 'L.01.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'gantry L4C', code: 'L.01.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'gantry E3', code: 'L.01.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'gantry E2', code: 'L.01.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'gantry E1', code: 'L.01.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L5', code: 'L.01.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L4', code: 'L.01.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'L.01.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'L.01.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '01-Floor between laser and experimental floors', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plenum', code: 'L.01.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'South corridor (for equipment)', code: 'L.02.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'North corridor (for personnel access)', code: 'L.02.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E1 experimental hall', code: 'L.02.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E1 small labyrinth', code: 'L.02.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E1 control room', code: 'L.02.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E2 control room', code: 'L.02.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E2 small labyrinth', code: 'L.02.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E2 experimental hall', code: 'L.02.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room', code: 'L.02.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room', code: 'L.02.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E3 small labyrinth', code: 'L.02.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E3 control room', code: 'L.02.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E3 experimental hall', code: 'L.02.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'laser technology L4c', code: 'L.02.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L1', code: 'L.02.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room', code: 'L.02.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'L.02.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E4 control room', code: 'L.02.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'L4c control room', code: 'L.02.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E4 small labyrinth', code: 'L.02.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E4 experimental hall', code: 'L.02.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L6', code: 'L.02.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room', code: 'L.02.24', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E6 experimental hall', code: 'L.02.25', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'L.02.26', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E6 control room', code: 'L.02.28', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'L.02.29', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L2', code: 'L.02.31', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E5 control room', code: 'L.02.32', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'L.02.33', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E5 experimental hall', code: 'L.02.34', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room', code: 'L.02.35', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'L.02.36', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L5', code: 'L.02.37', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room South (next to E5)', code: 'L.02.38', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L4', code: 'L.02.39', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'L.02.40', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'L.02.41', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'L.02.42', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'L.02.43', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E1 large labyrinth', code: 'L.02.44', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E2 large labyrinth', code: 'L.02.45', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E3 large labyrinth', code: 'L.02.46', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'E4 large labyrinth', code: 'L.02.47', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '02-Experimental floor', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'L.02.48', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '1-Floor between laser AND support floors', facility: 'B' })
MERGE (l3:Location{ name: "visitor's gallery", code: 'L.1.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '1-Floor between laser and support floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L1', code: 'L.1.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '1-Floor between laser and support floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L2', code: 'L.1.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '1-Floor between laser and support floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L3', code: 'L.1.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '1-Floor between laser and support floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L5', code: 'L.1.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '1-Floor between laser and support floors', facility: 'B' })
MERGE (l3:Location{ name: 'stair L4', code: 'L.1.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'L.2.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plantroom', code: 'L.2.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'North corridor (for personnel access)', code: 'L.2.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plantroom', code: 'L.2.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L1', code: 'L.2.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: "visitor's gallery top floor", code: 'L.2.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'L4a laser technology - storage room', code: 'L.2.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'South corridor (for equipment)', code: 'L.2.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'Liquid nitrogen plantroom', code: 'L2.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'technology - capacitors', code: 'L.2.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'support systems', code: 'L.2.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'L.2.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'L.2.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'L.2.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'liquid nitrogen plantroom', code: 'L.2.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L2', code: 'L.2.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair L3', code: 'L.2.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'L.2.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laser Building', facility: 'B' })
MERGE (l2:Location{ name: '2-Support technologies floor', facility: 'B' })
MERGE (l3:Location{ name: 'platform for vacuum pumps', code: 'L.3.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.00.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'stair LB1', code: 'LB.00.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'break out room staff', code: 'LB.00.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'meeting room 01', code: 'LB.00.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'meeting room 02', code: 'LB.00.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'office for technicians 01', code: 'LB.00.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'office for technicians 02', code: 'LB.00.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'office for technicians 03', code: 'LB.00.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'office for technicians 04', code: 'LB.00.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'security office', code: 'LB.00.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'LB.00.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'clean assembly room 01', code: 'LB.00.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'clean assembly room 02', code: 'LB.00.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'shower lobby', code: 'LB.00.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'shower', code: 'LB.00.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.00.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'changing room', code: 'LB.00.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'Main control room', code: 'LB.00.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'control room for oscillator', code: 'LB.00.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'control room for L4', code: 'LB.00.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'break out area', code: 'LB.00.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.00.24', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.00.25', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'first aid', code: 'LB.00.26', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'stair LB2', code: 'LB.00.27', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'changing room', code: 'LB.00.28', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'LB.00.29', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'large equipment airlock', code: 'LB.00.30', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.00.31', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.00.32', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'disabled access toilet', code: 'LB.00.33', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'LB.00.34', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'kitchenette', code: 'LB.00.37', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'LB.00.38', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'data switch room', code: 'LB.00.44', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'personal entrance', code: 'LB.00.45', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor (lasers)', facility: 'B' })
MERGE (l3:Location{ name: 'personal entrance', code: 'LB.00.46', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.01.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair LB1', code: 'LB.01.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'LB.01.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'optical workshop (CESNET reserved)', code: 'LB.01.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'optical workshop (CESNET reserved)', code: 'LB.01.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'UPS battery room ', code: 'LB.01.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'HV switch + transformer', code: 'LB.01.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'power switch room', code: 'LB.01.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Server room S2 ', code: 'LB.01.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.01.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'LB.01.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair LB2', code: 'LB.01.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'changing room', code: 'LB.01.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'material airlock', code: 'LB.01.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'ultrasonic cleaning', code: 'LB.01.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.01.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.01.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'LB.01.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (men)', code: 'LB.01.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (women)', code: 'LB.01.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.01.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'gas meter room', code: 'LB.01.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'fire power switch room', code: 'LB.01.27', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.01.29', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.01.30', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room', code: 'LB.01.31', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'changing room', code: 'LB.01.32', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'material air lock', code: 'LB.01.33', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'target lab', code: 'LB.01.34', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'LB.01.35', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'air lock Entrance', code: 'LB.01.36', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'air lock Exit', code: 'LB.01.37', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'LB.01.38', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'preparation', code: 'LB.01.39', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'cleaning room', code: 'LB.01.40', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '01- First undeground floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'LB.01.41', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'South corridor for large equipment', code: 'LB.02.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'stair LB1', code: 'LB.02.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'sprinkler plant room', code: 'LB.02.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'sprinkler tank room', code: 'LB.02.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'nanolab', code: 'LB.02.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'nanolab clean part', code: 'LB.02.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'air shower', code: 'LB.02.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'changing room', code: 'LB.02.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'LB.02.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'metrology lab', code: 'LB.02.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.02.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'assembly area', code: 'LB.02.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'assembly area', code: 'LB.02.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'changing room', code: 'LB.02.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'cleanroom attendant', code: 'LB.02.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.02.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'ultrasonic cleaning area', code: 'LB.02.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'large equipment airlock', code: 'LB.02.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.02.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.02.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'LB.02.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (men)', code: 'LB.02.24', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (women)', code: 'LB.02.25', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.02.26', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'LB.02.27', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'pump room', code: 'LB.02.28', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'LB.02.32', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'air shower', code: 'LB.02.33', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'air shower', code: 'LB.02.34', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'dozimetry room', code: 'LB.02.35', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'changing room', code: 'LB.02.36', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'material air lock', code: 'LB.02.37', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.02.38', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'dark room FM', code: 'LB.02.39', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'dark room ED', code: 'LB.02.40', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'biological laboratory', code: 'LB.02.41', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'cold room', code: 'LB.02.42', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cabinet", code: 'LB.02.42', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'chemical laboratory', code: 'LB.02.43', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: '1st Ante-room (clean)', code: 'LB.02.44', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: '2nd Ante-room (dirty)', code: 'LB.02.47', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'air lock', code: 'LB.02.48', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'plant room', code: 'LB.02.49', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'biological laboratory', code: 'LB.02.50', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'meeting room / lab manager', code: 'LB.02.51', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'material Air lock - dirty part', code: 'LB.02.52', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'material Air lock - clean part', code: 'LB.02.53', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby steam sterilisator', code: 'LB.02.54', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.02.55', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'biological laboratory', code: 'LB.02.56', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'main entrance changing room', code: 'LB.02.57', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cabinet", code: 'LB.02.59', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '02-Second underground floor (experiments)', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.02.61', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.1.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair LB1', code: 'LB.1.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'fine workshop 01', code: 'LB.1.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'fine workshop 02', code: 'LB.1.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'fine workshop 03', code: 'LB.1.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'fine workshop 04', code: 'LB.1.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'fine workshop 05', code: 'LB.1.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'fine workshop 06', code: 'LB.1.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'fine workshop 07', code: 'LB.1.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.1.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'LB.1.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'welding room', code: 'LB.1.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'shower lobby', code: 'LB.1.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'shower', code: 'LB.1.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.1.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'WC lobby (men)', code: 'LB.1.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'WC lobby (women)', code: 'LB.1.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.1.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'mechanical workshop', code: 'LB.1.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'workshop supervisor', code: 'LB.1.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.1.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'large equipment loading', code: 'LB.1.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.1.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'changing room (women)', code: 'LB.1.24', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'changing room (men)', code: 'LB.1.25', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'WC lobby (men)', code: 'LB.1.26', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'WC lobby (women)', code: 'LB.1.27', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'kitchenette', code: 'LB.1.28', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'personal entrance', code: 'LB.1.36', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'personal entrance', code: 'LB.1.37', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'data switch room', code: 'LB.1.38', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '1-first floor', facility: 'B' })
MERGE (l3:Location{ name: 'security point', code: 'LB.1.39', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.2.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair LB1', code: 'LB.2.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'LB.2.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'LB.2.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'LB.2.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'LB.2.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'LB.2.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'LB.2.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'LB.2.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.2.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'shower lobby', code: 'LB.2.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'shower', code: 'LB.2.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.2.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plantroom', code: 'LB.2.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'calibration laboratory', code: 'LB.2.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'changing room', code: 'LB.2.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'large equipment lobby', code: 'LB.2.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.2.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'LB.2.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.2.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'disabled access toilet', code: 'LB.2.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'LB.2.24', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (men)', code: 'LB.2.25', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (women)', code: 'LB.2.26', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'kitchenette', code: 'LB.2.27', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'personal entrance', code: 'LB.2.34', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'personal entrance', code: 'LB.2.35', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'data switch room', code: 'LB.2.36', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'airlock', code: 'LB.2.37', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'target laboratory', code: 'LB.2.38', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '3-third floor (plantrooms)', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'LB.3.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '3-third floor (plantrooms)', facility: 'B' })
MERGE (l3:Location{ name: 'stair LB1', code: 'LB.3.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '3-third floor (plantrooms)', facility: 'B' })
MERGE (l3:Location{ name: 'boiler + HC distribution', code: 'LB.3.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: '3-third floor (plantrooms)', facility: 'B' })
MERGE (l3:Location{ name: 'goods lift plantroom', code: 'LB.3.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: 'personal-goods lift', facility: 'B' })
MERGE (l3:Location{ name: '', code: 'LB.LI.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Laboratory building', facility: 'B' })
MERGE (l2:Location{ name: 'goods lift', facility: 'B' })
MERGE (l3:Location{ name: '', code: 'LB.LI.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'canteen', code: 'M.00.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (women)', code: 'M.00.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'M.00.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (men)', code: 'M.00.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'disabled access toilet', code: 'M.00.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's croom", code: 'M.00.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair M1', code: 'M.00.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (men)', code: 'M.00.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'wheelchair access. toilet', code: 'M.00.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (women)', code: 'M.00.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'refuse store', code: 'M.00.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'M.00.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'bldg management (2 p)', code: 'M.00.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'bldg management (2 p)', code: 'M.00.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'bldg management (2 p)', code: 'M.00.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'M.00.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'M.00.32', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'canteen and library', code: 'M.1.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor + break room', code: 'M.1.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair M2', code: 'M.1.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meeting room A', code: 'M.1.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meeting room C', code: 'M.1.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meeting room B', code: 'M.1.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'M.1.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'disabled access toilet', code: 'M.1.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (men)', code: 'M.1.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (women)', code: 'M.1.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair M1', code: 'M.1.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meeting room (dark room)', code: 'M.1.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'M.1.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'M.1.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'M.1.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'M.1.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'AV room', code: 'M.1.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'M.1.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Main lecture hall (150 p)', code: 'M.1.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets (men)', code: 'M.1.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets (women)', code: 'M.1.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'library', code: 'M.2.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'meeting room', code: 'M.2.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plantroom', code: 'M.2.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'M.2.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair M1', code: 'M.2.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'M.2.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1p)', code: 'M.2.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'M.2.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'M.2.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'entrance lobby', code: 'O.00.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.00.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'day room lobby', code: 'O.00.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'security day room', code: 'O.00.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'security shower', code: 'O.00.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Server room S1', code: 'O.00.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'security (3 p)', code: 'O.00.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'security (2 p)', code: 'O.00.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'entrance corridor', code: 'O.00.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'IT support (2 p)', code: 'O.00.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.00.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.00.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.00.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.00.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair O2', code: 'O.00.24', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'kitchenette', code: 'O.00.25', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'break out area', code: 'O.00.26', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (men)', code: 'O.00.27', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'disabled access toilet', code: 'O.00.28', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meeting room F', code: 'O.00.29', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'break out area', code: 'O.00.30', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'O.00.31', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meeting room G', code: 'O.00.32', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'O.00.33', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (women)', code: 'O.00.34', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'shower', code: 'O.00.35', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'kitchenette, dining', code: 'O.00.36', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'break out area', code: 'O.00.37', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair O1', code: 'O.00.38', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'O.00.39', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (4 p)', code: 'O.00.40', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'research laboratory (2 p)', code: 'O.00.41', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.00.42', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.00.43', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: '8 person open office', code: 'O.00.44', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: '8 person open office', code: 'O.00.45', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.00.46', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.00.47', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.00.48', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.00.49', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.00.50', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.00.51', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.00.52', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.00.53', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: '9 person open office', code: 'O.00.54', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'O.00.57', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'O.00.58', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'data switch room', code: 'O.00.59', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'technology control room', code: 'O.00.60', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '00-Ground floor', facility: 'B' })
MERGE (l3:Location{ name: 'technology room storage', code: 'O.00.61', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'O.1.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (6 p)', code: 'O.1.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.1.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'research laboratory (2 p)', code: 'O.1.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'research laboratory (2 p)', code: 'O.1.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'data switch room', code: 'O.1.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.1.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.1.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.1.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (4 p)', code: 'O.1.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair O2', code: 'O.1.24', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'kitchenette', code: 'O.1.25', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'O.1.26', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (men)', code: 'O.1.27', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'disabled access toilet', code: 'O.1.28', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meeting room H', code: 'O.1.29', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'O.1.30', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'O.1.31', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meting room I', code: 'O.1.32', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'O.1.33', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (women)', code: 'O.1.34', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'shower', code: 'O.1.35', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'kitchenette + dining', code: 'O.1.36', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'O.1.37', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair O1', code: 'O.1.38', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'O.1.39', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (6 p)', code: 'O.1.40', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (4 p)', code: 'O.1.41', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (4 p)', code: 'O.1.42', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.1.43', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.1.44', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: '8 person open office', code: 'O.1.45', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: '8 person open office', code: 'O.1.46', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.1.47', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.1.48', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.1.49', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.1.50', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.1.51', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.1.52', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.1.53', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.1.54', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: '9 person open office', code: 'O.1.55', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets (men)', code: 'O.1.56', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets (women)', code: 'O.1.57', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'O.1.58', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '1-First floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'O.1.59', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'lobby', code: 'O.2.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (6 p)', code: 'O.2.02', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.2.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.05', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.06', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.07', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.08', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.09', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.10', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.11', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'research laboratory (2 p)', code: 'O.2.12', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'research laboratory (2 p)', code: 'O.2.13', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.2.14', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'data switch room', code: 'O.2.15', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.2.16', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.17', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.18', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.19', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.20', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (2 p)', code: 'O.2.21', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.2.22', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (4 p)', code: 'O.2.23', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair O2', code: 'O.2.24', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'kitchenette', code: 'O.2.25', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'O.2.26', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meeting room J', code: 'O.2.29', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'O.2.30', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'O.2.31', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'Meeting room K', code: 'O.2.32', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'O.2.33', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'toilets lobby (women)', code: 'O.2.34', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'kitchenette + dining', code: 'O.2.36', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'corridor', code: 'O.2.37', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'stair O1', code: 'O.2.38', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: "cleaner's cupboard", code: 'O.2.39', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (6 p)', code: 'O.2.40', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (4 p)', code: 'O.2.41', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.2.42', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.2.43', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.2.44', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: '8 person open office', code: 'O.2.45', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: '8 person open office', code: 'O.2.46', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.2.47', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.2.48', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.2.49', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.2.50', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (3 p)', code: 'O.2.51', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.2.52', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.2.53', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'office (1 p)', code: 'O.2.54', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: '9 person open office', code: 'O.2.55', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'O.2.58', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'storage', code: 'O.2.59', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: '2-Second floor', facility: 'B' })
MERGE (l3:Location{ name: 'HVAC plant room', code: 'O.3.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (l2)-[:HAS_SUBLOCATION]->(l3)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: 'Personal lift', code: 'O.LI.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Office building', facility: 'B' })
MERGE (l2:Location{ name: 'Personal lift', code: 'O.LI.04', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Technical gases building', facility: 'B' })
MERGE (l2:Location{ name: 'tech. gases plantroom', code: 'TG.00.01', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);

MATCH (f:Facility{ code: 'B' })
MERGE (l1:Location{ name: 'Multifunctional building', facility: 'B' })
MERGE (l2:Location{ name: ' Personal Goods Lift', code: 'M.LI.03', facility: 'B' })
MERGE (l1)-[:HAS_SUBLOCATION]->(l2)
MERGE (f)-[:HAS_LOCATION]->(l1);
