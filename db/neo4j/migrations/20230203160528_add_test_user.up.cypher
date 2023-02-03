MATCH (basics:Role{ uid:'6b535e07-2a1c-4498-a1b3-8d63f9275519' })

CREATE (testUser:User{
  uid: 'ec9ba859-219a-4067-a95e-72580d2fa787',
  username: 'test',
  firstName: 'Test',
  lastName: 'Tester',
  email: 'panda.test@eli-beams.eu',
  passwordHash: '$2a$12$.tOnYIBKOBevvISuttNiS..cjfrag2b9PljB4gtp9GRyNn5Kp5Wy.',
  isEnabled: true })
  
  CREATE(testUser)-[:HAS_ROLE]->(basics)
