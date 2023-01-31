// Create admin user
MERGE (admin:User{
  uid: '71864520-9e86-427c-901c-0c220f951775',
  username: 'admin',
  firstName: 'Jiří',
  lastName: 'Švácha',
  email: 'jiri.svacha@eli-beams.eu',
  passwordHash: '$2a$13$ifCe51bH2rvTlAH2F1DLnuyCdM.yHt.KhAASXXhjQeGVHYwy3RdSO',
  isEnabled: true })
