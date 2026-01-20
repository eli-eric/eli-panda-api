MERGE (r:Researcher {identificationNumber: "840404"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Smijesh", r.lastName = "Achary", r.identificationNumber = "840404", r.orcid = "0000-0002-4884-3336", r.scopusId = "35265249300", r.researcherId = "E-9444-2015"
WITH r
MATCH (c:Country {code: "IN"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9101123427"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martin", r.lastName = "Albrecht", r.identificationNumber = "9101123427", r.orcid = "0000-0002-0103-5614", r.scopusId = "57020873300", r.researcherId = "GZZ-7943-2022"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "750328"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jakob", r.lastName = "Andreasson", r.identificationNumber = "750328", r.orcid = "0000-0002-3202-2330"
WITH r
MATCH (c:Country {code: "SE"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "6706062132"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Viktor", r.lastName = "Andrle", r.identificationNumber = "6706062132"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "670915"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Borislav", r.lastName = "Angelov", r.identificationNumber = "670915", r.orcid = "0000-0003-1667-0626", r.scopusId = "6603229727"
WITH r
MATCH (c:Country {code: "BG"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "820731"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Roman", r.lastName = "Antipenkov", r.identificationNumber = "820731", r.orcid = "0000-0003-0565-1710"
WITH r
MATCH (c:Country {code: "LT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "810205"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Anabella", r.lastName = "Araudo", r.identificationNumber = "810205", r.orcid = "0000-0001-7605-5786", r.scopusId = "15077922600"
WITH r
MATCH (c:Country {code: "AR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "6553152353"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Libuše", r.lastName = "Arnoštová", r.identificationNumber = "6553152353"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "6805100434"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Pavel", r.lastName = "Bakule", r.identificationNumber = "6805100434", r.orcid = "0000-0002-7417-5810", r.scopusId = "6603364608"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "940531"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Iullia", r.lastName = "Baranova", r.identificationNumber = "940531", r.orcid = "0009-0003-2262-2054"
WITH r
MATCH (c:Country {code: "RU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8806033489"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jan", r.lastName = "Bartoníček", r.identificationNumber = "8806033489"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7609273045"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Radek", r.lastName = "Baše", r.identificationNumber = "7609273045"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "760205"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "César", r.lastName = "Bernardo", r.identificationNumber = "760205", r.orcid = "0000-0002-7372-083X", r.scopusId = "56102755500"
WITH r
MATCH (c:Country {code: "PT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8701704605"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Robert", r.lastName = "Boge", r.identificationNumber = "8701704605", r.orcid = "0000-0002-8240-0854"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8805115803"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Karel", r.lastName = "Boháček", r.identificationNumber = "8805115803", r.scopusId = "56593111300", r.researcherId = "HKH-1227-2023"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "470816"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Sergey", r.lastName = "Bulanov", r.identificationNumber = "470816", r.orcid = "0000-0001-8305-0289"
WITH r
MATCH (c:Country {code: "JP"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "811014"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Anna", r.lastName = "Cimmino", r.identificationNumber = "811014", r.orcid = "0000-0001-7510-4996", r.researcherId = "AAA-1673-2020"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "740820"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Pablo", r.lastName = "Cirrone", r.identificationNumber = "740820", r.orcid = "0000-0001-5733-9281", r.scopusId = "6604038461", r.researcherId = "I-6474-2015"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "910213"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Florian", r.lastName = "Condamine", r.identificationNumber = "910213", r.orcid = "0000-0002-4163-6895", r.scopusId = "57190088853"
WITH r
MATCH (c:Country {code: "FR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8910121319"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martin", r.lastName = "Cuhra", r.identificationNumber = "8910121319", r.orcid = "0000-0002-3348-7914"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8906160131"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Josef", r.lastName = "Cupal", r.identificationNumber = "8906160131", r.orcid = "0000-0003-2609-4597"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9006262155"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jan", r.lastName = "Černý", r.identificationNumber = "9006262155"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9262178915"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Petra", r.lastName = "Čubáková", r.identificationNumber = "9262178915", r.orcid = "0000-0003-4118-4662"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8311195321"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jakub", r.lastName = "Dostál", r.identificationNumber = "8311195321", r.orcid = "0000-0003-2188-0863", r.scopusId = "55316098600", r.researcherId = "X-6931-2019"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9354184730"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Petra", r.lastName = "Dvořáková Ruskayová", r.identificationNumber = "9354184730", r.orcid = "0000-0001-8626-1292", r.scopusId = "57221596433"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9109112639"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jan", r.lastName = "Eisenschreiber", r.identificationNumber = "9109112639", r.orcid = "0000-0003-3295-9120"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "920728"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Emily", r.lastName = "Erdman", r.identificationNumber = "920728", r.orcid = "0000-0002-3946-6436"
WITH r
MATCH (c:Country {code: "US"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "826223"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Shirly Josefina", r.lastName = "Espinoza Herrera", r.identificationNumber = "826223", r.orcid = "0000-0002-2740-7156"
WITH r
MATCH (c:Country {code: "AR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "990809"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Gaëtan", r.lastName = "Fauvel", r.identificationNumber = "990809", r.orcid = "0000-0002-1088-3412"
WITH r
MATCH (c:Country {code: "FR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "900722"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Yana", r.lastName = "Fetisova", r.identificationNumber = "900722", r.orcid = "0000-0002-6711-1819"
WITH r
MATCH (c:Country {code: "RU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8106284131"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martin", r.lastName = "Fibrich", r.identificationNumber = "8106284131", r.orcid = "0000-0002-4121-0181", r.scopusId = "22934051700"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9309060046"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Ondřej", r.lastName = "Finke", r.identificationNumber = "9309060046", r.orcid = "0000-0002-3961-5446", r.scopusId = "57202090136", r.researcherId = "DWP-4359-2022"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8808114480"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martin", r.lastName = "Formánek", r.identificationNumber = "8808114480", r.orcid = "0000-0003-2704-6474", r.scopusId = "57103140600"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7005195087"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Tomáš", r.lastName = "Franek", r.identificationNumber = "7005195087"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8403313512"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jiří", r.lastName = "Freisleben", r.identificationNumber = "8403313512"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "860522"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Gavin", r.lastName = "Friedman", r.identificationNumber = "860522", r.orcid = "0000-0003-3462-309X", r.scopusId = "57195127533"
WITH r
MATCH (c:Country {code: "US"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "880130"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Nina", r.lastName = "Gamaiunova", r.identificationNumber = "880130", r.orcid = "0000-0002-8604-8245", r.scopusId = "56366177400"
WITH r
MATCH (c:Country {code: "UA"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "841124"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Evgeny", r.lastName = "Gelfer", r.identificationNumber = "841124", r.orcid = "0000-0002-6800-0483", r.scopusId = "27667642300"
WITH r
MATCH (c:Country {code: "RU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "920909"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Antoine", r.lastName = "Gintrand", r.identificationNumber = "920909", r.orcid = "0000-0001-8182-4324"
WITH r
MATCH (c:Country {code: "FR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "790411"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Lorenzo", r.lastName = "Giuffrida", r.identificationNumber = "790411", r.orcid = "0000-0003-4145-4829", r.scopusId = "16232797400"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7710125082"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jiří", r.lastName = "Golasowski", r.identificationNumber = "7710125082"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "690811"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Leonardo", r.lastName = "Goncalves", r.identificationNumber = "690811", r.scopusId = "57211259305"
WITH r
MATCH (c:Country {code: "PT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "801016"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jonathan Tyler", r.lastName = "Green", r.identificationNumber = "801016", r.orcid = "0000-0001-7027-5278"
WITH r
MATCH (c:Country {code: "US"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "980512"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Annika", r.lastName = "Grenfell", r.identificationNumber = "980512", r.orcid = "0000-0002-3314-4639"
WITH r
MATCH (c:Country {code: "FI"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9103205221"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Filip", r.lastName = "Grepl", r.identificationNumber = "9103205221", r.orcid = "0000-0003-4959-3575", r.scopusId = "57203638439"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9151306164"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martina", r.lastName = "Greplová Žáková", r.identificationNumber = "9151306164", r.orcid = "0000-0003-4352-5299"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "891121"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Gabriele Maria", r.lastName = "Grittani", r.identificationNumber = "891121", r.orcid = "0000-0003-4394-854X", r.scopusId = "55655471700", r.researcherId = "H-3158-2014"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "840827"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Yanjun", r.lastName = "Gu", r.identificationNumber = "840827", r.orcid = "0000-0002-6234-8489", r.scopusId = "36092692900", r.researcherId = "G-5767-2014"
WITH r
MATCH (c:Country {code: "CN"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "960803"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Arsenios", r.lastName = "Hadjikyriacou", r.identificationNumber = "960803", r.orcid = "0000-0002-0169-9074"
WITH r
MATCH (c:Country {code: "CY"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "860714"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Prokopis", r.lastName = "Hadjisolomou", r.identificationNumber = "860714", r.orcid = "0000-0003-1170-7397", r.scopusId = "57188864239", r.researcherId = "DXG-9326-2022"
WITH r
MATCH (c:Country {code: "CY"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8507106663"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Ľudovít", r.lastName = "Haizer", r.identificationNumber = "8507106663"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "480917"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Janos", r.lastName = "Hajdu", r.identificationNumber = "480917", r.orcid = "0000-0002-3747-2760"
WITH r
MATCH (c:Country {code: "HU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9452241469"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Irena", r.lastName = "Havlíčková", r.identificationNumber = "9452241469"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7908145344"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Mojmír", r.lastName = "Havlík", r.identificationNumber = "7908145344", r.orcid = "0000-0002-2434-1526"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "751227"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Juan Carlos", r.lastName = "Hernandez Martin", r.identificationNumber = "751227"
WITH r
MATCH (c:Country {code: "ES"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7802180067"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Pavel", r.lastName = "Homer", r.identificationNumber = "7802180067"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7808040009"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Aleš", r.lastName = "Honsa", r.identificationNumber = "7808040009"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "830215"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Ziaul", r.lastName = "Hoque", r.identificationNumber = "830215", r.orcid = "0000-0001-7869-6659"
WITH r
MATCH (c:Country {code: "BD"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8710153793"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martin", r.lastName = "Horáček", r.identificationNumber = "8710153793", r.orcid = "0000-0002-7176-1339"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8008202576"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Radek", r.lastName = "Horálek", r.identificationNumber = "8008202576"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8604071465"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Ondřej", r.lastName = "Hort", r.identificationNumber = "8604071465", r.orcid = "0000-0002-1330-0825", r.scopusId = "55115714400", r.researcherId = "IVS-6422-2023"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "840812"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Dávid", r.lastName = "Horváth", r.identificationNumber = "840812", r.orcid = "0000-0003-0986-9119"
WITH r
MATCH (c:Country {code: "HU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8304260063"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jan", r.lastName = "Hřebíček", r.identificationNumber = "8304260063"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "480114452"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Petr", r.lastName = "Hříbek", r.identificationNumber = "480114452", r.orcid = "0000-0003-3344-1083"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "760224"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Hsiao-Chih", r.lastName = "Huang", r.identificationNumber = "760224", r.orcid = "0000-0003-3428-266X"
WITH r
MATCH (c:Country {code: "TW"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8512151142"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jan", r.lastName = "Hubáček", r.identificationNumber = "8512151142"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "900317"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Andreas", r.lastName = "Hult Roos", r.identificationNumber = "900317", r.orcid = "0000-0001-8919-0979"
WITH r
MATCH (c:Country {code: "SE"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "711230"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Edwin", r.lastName = "Chacon Golcher", r.identificationNumber = "711230", r.orcid = "0000-0001-5322-425X"
WITH r
MATCH (c:Country {code: "CR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7907169974"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Timofej", r.lastName = "Chagovets", r.identificationNumber = "7907169974", r.orcid = "0000-0001-5943-3845", r.researcherId = "G-9079-2014"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "810607"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Uddhab", r.lastName = "Chaulagain", r.identificationNumber = "810607", r.orcid = "0000-0003-0125-0303", r.scopusId = "55607357900"
WITH r
MATCH (c:Country {code: "NP"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "900914"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Adrien", r.lastName = "Chauvin", r.identificationNumber = "900914", r.orcid = "0000-0003-1896-8776", r.scopusId = "56509186500"
WITH r
MATCH (c:Country {code: "FR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9005224481"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Lukáš", r.lastName = "Indra", r.identificationNumber = "9005224481", r.orcid = "0000-0002-2712-7242"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "950408"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Valeriia", r.lastName = "Istokskaia", r.identificationNumber = "950408", r.orcid = "0000-0003-4768-9207", r.scopusId = "57204909263"
WITH r
MATCH (c:Country {code: "RU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "721106"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Rachael", r.lastName = "Jack", r.identificationNumber = "721106", r.orcid = "0000-0002-9482-9316"
WITH r
MATCH (c:Country {code: "GB"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "5807231045"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Alexandr", r.lastName = "Jančárek", r.identificationNumber = "5807231045", r.orcid = "0000-0003-2510-6592", r.scopusId = "6505830458"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7912110272"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jakub", r.lastName = "Janďourek", r.identificationNumber = "7912110272"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "780804"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jeffrey Alan", r.lastName = "Jarboe", r.identificationNumber = "780804"
WITH r
MATCH (c:Country {code: "US"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8659099328"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Lucia", r.lastName = "Jarboe", r.identificationNumber = "8659099328", r.orcid = "0000-0003-2658-2929"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "700111"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Tae Moon", r.lastName = "Jeong", r.identificationNumber = "700111", r.orcid = "0000-0001-7206-2401"
WITH r
MATCH (c:Country {code: "KR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8711154771"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martin", r.lastName = "Jirka", r.identificationNumber = "8711154771", r.orcid = "0000-0003-4457-4471", r.scopusId = "55807086700", r.researcherId = "H-4152-2014"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "920610"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Noémie", r.lastName = "Jourdain", r.identificationNumber = "920610", r.orcid = "0000-0001-7793-2816", r.scopusId = "57201196754", r.researcherId = "DYW-0979-2022"
WITH r
MATCH (c:Country {code: "FR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9604154890"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Lucie", r.lastName = "Jurkovičová", r.identificationNumber = "9604154890", r.orcid = "0000-0002-4454-8268", r.scopusId = "57211499707", r.researcherId = "JYM-9767-2024"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8253270179"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Hedvika", r.lastName = "Kadlecová", r.identificationNumber = "8253270179", r.orcid = "0000-0003-4066-2175", r.scopusId = "23060832700", r.researcherId = "GCD-7795-2022"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "710326"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Vasiliki", r.lastName = "Kantarelou", r.identificationNumber = "710326", r.orcid = "0000-0003-0341-0016", r.scopusId = "30967548000"
WITH r
MATCH (c:Country {code: "GR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "5901040497"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jaroslav", r.lastName = "Kašpar", r.identificationNumber = "5901040497"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "870717"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Krishna Prasad", r.lastName = "Khakurel", r.identificationNumber = "870717", r.orcid = "0000-0001-5974-8302", r.scopusId = "57112467200"
WITH r
MATCH (c:Country {code: "NP"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "890323"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Alex", r.lastName = "Kirby", r.identificationNumber = "890323", r.orcid = "0000-0002-5597-5564"
WITH r
MATCH (c:Country {code: "GB"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8362117676"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Eva", r.lastName = "Klimešová", r.identificationNumber = "8362117676", r.orcid = "0000-0002-9569-7511", r.scopusId = "55376496100"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8004210049"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Ondřej", r.lastName = "Klimo", r.identificationNumber = "8004210049", r.orcid = "0000-0002-0565-2409", r.scopusId = "6506792992", r.researcherId = "B-2196-2010"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8411042519"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Miroslav", r.lastName = "Kloz", r.identificationNumber = "8411042519", r.orcid = "0000-0003-4609-8018", r.scopusId = "37119035600", r.researcherId = "H-3074-2014"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "6306186326"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Viliam", r.lastName = "Kmetík", r.identificationNumber = "6306186326", r.orcid = "0000-0003-1846-1975"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "521021"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Georg", r.lastName = "Korn", r.identificationNumber = "521021", r.orcid = "0000-0002-7093-5296"
WITH r
MATCH (c:Country {code: "DE"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "6105252021"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Michal", r.lastName = "Košelja", r.identificationNumber = "6105252021"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8610032013"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Pavel", r.lastName = "Koupil", r.identificationNumber = "8610032013", r.orcid = "0000-0003-3332-3503", r.scopusId = "57440016000", r.researcherId = "AAB-9565-2022"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7361180035"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Michaela", r.lastName = "Kozlová", r.identificationNumber = "7361180035", r.orcid = "0000-0003-4757-9972"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7909041591"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Daniel", r.lastName = "Kramer", r.identificationNumber = "7909041591", r.orcid = "0000-0003-3885-9198", r.scopusId = "7203031817", r.researcherId = "D-9840-2011"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "760327"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Maria", r.lastName = "Krikunova", r.identificationNumber = "760327", r.orcid = "0000-0002-6152-1825", r.scopusId = "7801545033", r.researcherId = "M-6805-2017"
WITH r
MATCH (c:Country {code: "DE"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "461123"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jean Claude", r.lastName = "Lagron", r.identificationNumber = "461123"
WITH r
MATCH (c:Country {code: "FR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9306303853"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Marcel", r.lastName = "Lamač", r.identificationNumber = "9306303853", r.orcid = "0000-0001-9218-6073", r.scopusId = "57211499673", r.researcherId = "GBC-3777-2022"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7410133258"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Tomáš", r.lastName = "Laštovička", r.identificationNumber = "7410133258", r.orcid = "0000-0001-5636-5197"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "890517"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Carlo Maria", r.lastName = "Lazzarini", r.identificationNumber = "890517", r.orcid = "0000-0003-0750-8272"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "900722"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Benoit", r.lastName = "Lefebvre", r.identificationNumber = "900722", r.orcid = "0000-0001-8212-6624", r.scopusId = "57143949200"
WITH r
MATCH (c:Country {code: "CA"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "850426"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Nils", r.lastName = "Lenngren", r.identificationNumber = "850426", r.orcid = "0000-0001-7563-9843", r.scopusId = "55521120600"
WITH r
MATCH (c:Country {code: "SE"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "860216"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Yingliang", r.lastName = "Liu", r.identificationNumber = "860216", r.orcid = "0000-0002-4089-688X", r.scopusId = "57215930419", r.researcherId = "M-7865-2017"
WITH r
MATCH (c:Country {code: "CN"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9412134204"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Sebastian", r.lastName = "Lorenz", r.identificationNumber = "9412134204", r.orcid = "0000-0002-1620-7976"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "5502280146"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Ivan", r.lastName = "Lukačina", r.identificationNumber = "5502280146"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "930421"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Alexander John", r.lastName = "Macleod", r.identificationNumber = "930421", r.orcid = "0000-0002-0800-4423", r.scopusId = "57219621693", r.researcherId = "GBZ-8065-2022"
WITH r
MATCH (c:Country {code: "GB"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "930309"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Srimanta", r.lastName = "Maity", r.identificationNumber = "930309", r.orcid = "0000-0002-3841-444X", r.scopusId = "57201485737", r.researcherId = "JSB-3109-2023"
WITH r
MATCH (c:Country {code: "IN"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8907050614"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Karel", r.lastName = "Majer", r.identificationNumber = "8907050614", r.orcid = "0000-0002-5318-2528", r.researcherId = "Q-9578-2017"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "801021"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Daniele", r.lastName = "Margarone", r.identificationNumber = "801021", r.orcid = "0000-0002-1917-9683", r.researcherId = "HPH-0722-2023"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9105023620"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martin", r.lastName = "Matys", r.identificationNumber = "9105023620", r.orcid = "0000-0002-6215-4246", r.scopusId = "57200922424", r.researcherId = "AAG-3113-2020"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7905313504"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Tomáš", r.lastName = "Mazanec", r.identificationNumber = "7905313504", r.orcid = "0000-0002-6940-6143"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9002102736"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Petr", r.lastName = "Mazůrek", r.identificationNumber = "9002102736"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "590403"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Alexander", r.lastName = "Molodozhentsev", r.identificationNumber = "590403", r.orcid = "0000-0001-7633-761X", r.scopusId = "6603646494", r.researcherId = "FNW-5027-2022"
WITH r
MATCH (c:Country {code: "RU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "890501"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Alamgir", r.lastName = "Mondal", r.identificationNumber = "890501", r.orcid = "0000-0001-9908-1354", r.scopusId = "57200449788"
WITH r
MATCH (c:Country {code: "IN"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "920520"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Quentin", r.lastName = "Moreno-Gelos", r.identificationNumber = "920520", r.orcid = "0000-0002-6331-1637", r.scopusId = "57202038376"
WITH r
MATCH (c:Country {code: "FR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "920617"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Atripan", r.lastName = "Mukherjee", r.identificationNumber = "920617", r.orcid = "0000-0003-2303-9007", r.scopusId = "57216591968"
WITH r
MATCH (c:Country {code: "IN"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "910107"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Ana Laura", r.lastName = "Müller", r.identificationNumber = "910107", r.orcid = "0000-0002-8473-695X"
WITH r
MATCH (c:Country {code: "AR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8307141997"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jaroslav", r.lastName = "Nejdl", r.identificationNumber = "8307141997", r.orcid = "0000-0003-0864-8592", r.scopusId = "23036374800", r.researcherId = "G-5995-2014"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {firstName: "Jaromír", lastName: "Němeček"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jaromír", r.lastName = "Němeček", r.orcid = "0000-0002-6215-4246", r.scopusId = "57200922424"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "711124"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Gábor Balázs", r.lastName = "Németh", r.identificationNumber = "711124"
WITH r
MATCH (c:Country {code: "HU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8205270106"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Michal", r.lastName = "Nevrkla", r.identificationNumber = "8205270106", r.orcid = "0000-0002-1881-7882", r.scopusId = "38761986400"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "950811"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Sebastian", r.lastName = "Niekrasz", r.identificationNumber = "950811", r.orcid = "0000-0002-0428-5192"
WITH r
MATCH (c:Country {code: "PL"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8412122565"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jakub", r.lastName = "Novák", r.identificationNumber = "8412122565", r.orcid = "0000-0003-4903-2479", r.researcherId = "G-6891-2014"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7851250550"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Veronika", r.lastName = "Olšovcová", r.identificationNumber = "7851250550", r.orcid = "0000-0002-7955-9231", r.scopusId = "56617519000"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "6404021360"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Václav", r.lastName = "Orna", r.identificationNumber = "6404021360"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "980903"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Paul", r.lastName = "Pandikian", r.identificationNumber = "980903", r.scopusId = "57995180500", r.researcherId = "HIA-4471-2022"
WITH r
MATCH (c:Country {code: "FR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8901262843"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Tomáš", r.lastName = "Parkman", r.identificationNumber = "8901262843", r.orcid = "0000-0002-0425-1291", r.scopusId = "56989869400"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "740324"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Davorin", r.lastName = "Peceli", r.identificationNumber = "740324", r.orcid = "0000-0002-6672-0786"
WITH r
MATCH (c:Country {code: "HR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "910403"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Giada", r.lastName = "Petringa", r.identificationNumber = "910403", r.orcid = "0000-0002-2854-8436", r.scopusId = "56563089600", r.researcherId = "DYG-4875-2022"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "841216"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Alessandra", r.lastName = "Picchiotti", r.identificationNumber = "841216", r.orcid = "0000-0003-0167-1431", r.scopusId = "56341059900", r.researcherId = "DLG-9321-2022"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "890807"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Birgit", r.lastName = "Plötzeneder", r.identificationNumber = "890807", r.orcid = "0009-0000-0199-6522"
WITH r
MATCH (c:Country {code: "DE"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8903105739"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Petr", r.lastName = "Pokorný", r.identificationNumber = "8903105739"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "880319"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Vitaly", r.lastName = "Polovinkin", r.identificationNumber = "880319", r.orcid = "0000-0002-3630-5565", r.researcherId = "G-6627-2017"
WITH r
MATCH (c:Country {code: "RU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8112230038"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martin", r.lastName = "Přeček", r.identificationNumber = "8112230038", r.orcid = "0000-0002-5790-5543", r.scopusId = "35741533300", r.researcherId = "G-5648-2014"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8203123324"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jan", r.lastName = "Pšikal", r.identificationNumber = "8203123324", r.orcid = "0000-0003-4586-1149", r.scopusId = "15056568000", r.researcherId = "G-8403-2014"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "961027"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Yelyzaveta", r.lastName = "Pulnova", r.identificationNumber = "961027", r.orcid = "0000-0001-6562-5499", r.scopusId = "57222046139"
WITH r
MATCH (c:Country {code: "UA"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "530305098"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Ladislav", r.lastName = "Půst", r.identificationNumber = "530305098", r.orcid = "0000-0002-3471-9342", r.researcherId = "D-9600-2019"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9110104784"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Marek", r.lastName = "Raclavský", r.identificationNumber = "9110104784", r.orcid = "0000-0003-4836-9934", r.scopusId = "56389565100", r.researcherId = "DOD-1002-2022"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "770901"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Mateusz", r.lastName = "Rebarz", r.identificationNumber = "770901", r.orcid = "0000-0002-5823-2432", r.researcherId = "F-7739-2015"
WITH r
MATCH (c:Country {code: "PL"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "440102036"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Oldřich", r.lastName = "Renner", r.identificationNumber = "440102036", r.orcid = "0000-0003-4942-2637", r.scopusId = "7004055743", r.researcherId = "C-1591-2010"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "900317"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Andreas Hult ", r.lastName = "Roos", r.identificationNumber = "900317", r.orcid = "0000-0001-8919-0979"
WITH r
MATCH (c:Country {code: "SE"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8610138999"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Peter", r.lastName = "Rubovič", r.identificationNumber = "8610138999", r.scopusId = "54401686600"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "6308291594"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Bedřich", r.lastName = "Rus", r.identificationNumber = "6308291594", r.orcid = "0000-0002-8982-1029"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "511118168"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jan", r.lastName = "Řídký", r.identificationNumber = "511118168", r.orcid = "0000-0001-6697-1393", r.scopusId = "55944260100", r.researcherId = "H-6184-2014"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "470427"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Pavel", r.lastName = "Sasorov", r.identificationNumber = "470427", r.orcid = "0000-0001-9257-6565", r.scopusId = "7004127648"
WITH r
MATCH (c:Country {code: "RU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "810212"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Valentina", r.lastName = "Scuderi", r.identificationNumber = "810212", r.orcid = "0000-0003-2547-8978", r.scopusId = "23095874300", r.researcherId = "DQL-5354-2022"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9051283263"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jitka", r.lastName = "Sedmidubská", r.identificationNumber = "9051283263"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "610307"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Rashid", r.lastName = "Shaisultanov", r.identificationNumber = "610307", r.scopusId = "6507848862", r.researcherId = "DQF-6786-2022"
WITH r
MATCH (c:Country {code: "RU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "921019"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Sviatoslav", r.lastName = "Shekhanov", r.identificationNumber = "921019", r.orcid = "0000-0002-2125-8962", r.researcherId = "HKW-6944-2023"
WITH r
MATCH (c:Country {code: "RU"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "790606"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Francesco", r.lastName = "Schillaci", r.identificationNumber = "790606", r.orcid = "0000-0003-3628-5880", r.scopusId = "55768035400"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "870116"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Raj Laxmi", r.lastName = "Singh", r.identificationNumber = "870116", r.orcid = "0000-0001-7383-9753", r.scopusId = "56517259400"
WITH r
MATCH (c:Country {code: "IN"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "970624"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Anastasia", r.lastName = "Sklia", r.identificationNumber = "970624"
WITH r
MATCH (c:Country {code: "GR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9756215920"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Vanda", r.lastName = "Sluková", r.identificationNumber = "9756215920", r.orcid = "0000-0002-8594-2173"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7905083175"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "David", r.lastName = "Snopek", r.identificationNumber = "7905083175"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7605217675"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Stanislav", r.lastName = "Stanček", r.identificationNumber = "7605217675", r.orcid = "0000-0002-1695-2458"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9701202720"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Matyáš", r.lastName = "Staněk", r.identificationNumber = "9701202720", r.orcid = "0009-0002-1432-1721"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9301073089"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Vojtěch", r.lastName = "Stránský", r.identificationNumber = "9301073089", r.orcid = "0000-0001-5938-5118", r.scopusId = "57202765399"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8605185435"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Petr", r.lastName = "Szotkowski", r.identificationNumber = "8605185435", r.orcid = "0000-0002-0124-2517"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "920410"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Wojciech Jerzy", r.lastName = "Szuba", r.identificationNumber = "920410", r.orcid = "0000-0003-2133-2939"
WITH r
MATCH (c:Country {code: "PL"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9408290683"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Michal", r.lastName = "Šesták", r.identificationNumber = "9408290683"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9606226157"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jiří", r.lastName = "Šišma", r.identificationNumber = "9606226157", r.orcid = "0000-0001-7704-136X"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8703232120"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Václav", r.lastName = "Šobr", r.identificationNumber = "8703232120", r.orcid = "0000-0003-4181-5225"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9412292725"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Alexandr", r.lastName = "Špaček", r.identificationNumber = "9412292725", r.orcid = "0000-0002-9566-4699"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "480130"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Vladimir", r.lastName = "Tikhonchuk", r.identificationNumber = "480130", r.orcid = "0000-0001-7532-5879", r.researcherId = "S-1160-2018"
WITH r
MATCH (c:Country {code: "FR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "870321"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Tomas", r.lastName = "Tolenis", r.identificationNumber = "870321", r.orcid = "0000-0001-5400-729X"
WITH r
MATCH (c:Country {code: "LT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "870403"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Murat", r.lastName = "Torun", r.identificationNumber = "870403", r.orcid = "0000-0002-0534-6946", r.scopusId = "57193489208", r.researcherId = "I-7724-2016"
WITH r
MATCH (c:Country {code: "TR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "951026"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Marco", r.lastName = "Tosca", r.identificationNumber = "951026", r.orcid = "0000-0002-5330-3123", r.scopusId = "57219146450"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "7908255421"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Pavel", r.lastName = "Trojek", r.identificationNumber = "7908255421", r.orcid = "0000-0001-7841-4893", r.scopusId = "6507728520"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "6402161304"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Roman", r.lastName = "Truneček", r.identificationNumber = "6402161304"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "6455032133"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Zuzana", r.lastName = "Trunečková", r.identificationNumber = "6455032133"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "890905"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Maksym", r.lastName = "Tryus", r.identificationNumber = "890905", r.orcid = "0000-0003-0825-5051"
WITH r
MATCH (c:Country {code: "UA"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "850701"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Özlem", r.lastName = "Tunc", r.identificationNumber = "850701"
WITH r
MATCH (c:Country {code: "TR"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "830318"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Boguslaw", r.lastName = "Tykalewicz", r.identificationNumber = "830318"
WITH r
MATCH (c:Country {code: "PL"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9109285229"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Jan", r.lastName = "Vábek", r.identificationNumber = "9109285229", r.orcid = "0000-0001-7689-7174"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9204223578"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Petr", r.lastName = "Valenta", r.identificationNumber = "9204223578", r.orcid = "0000-0003-1067-3762", r.scopusId = "57204242912", r.researcherId = "AAN-6606-2021"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "830328"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Saul", r.lastName = "Vazquez Miranda", r.identificationNumber = "830328", r.orcid = "0000-0001-9897-5508"
WITH r
MATCH (c:Country {code: "MX"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "761229"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Andriy", r.lastName = "Velyhan", r.identificationNumber = "761229", r.orcid = "0000-0003-0074-3632", r.scopusId = "55910401100"
WITH r
MATCH (c:Country {code: "UA"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "750716"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Roberto", r.lastName = "Versaci", r.identificationNumber = "750716", r.orcid = "0000-0002-6516-0764", r.scopusId = "6603621883", r.researcherId = "AAG-6712-2019"
WITH r
MATCH (c:Country {code: "IT"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9511213646"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Filip", r.lastName = "Vitha", r.identificationNumber = "9511213646", r.orcid = "0000-0002-4350-5485"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8955310023"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Alžběta", r.lastName = "Vosátková", r.identificationNumber = "8955310023"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8902260191"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Štěpán", r.lastName = "Vyhlídka", r.identificationNumber = "8902260191", r.orcid = "0000-0002-1660-2175"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "640206"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Stefan Andreas", r.lastName = "Weber", r.identificationNumber = "640206", r.orcid = "0000-0003-3154-9306", r.scopusId = "7401926973"
WITH r
MATCH (c:Country {code: "DE"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "761117"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Tuomas", r.lastName = "Wiste", r.identificationNumber = "761117", r.scopusId = "35108286700", r.researcherId = "EFW-7118-2022"
WITH r
MATCH (c:Country {code: "FI"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "9004161463"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Martin", r.lastName = "Zahradník", r.identificationNumber = "9004161463", r.orcid = "0000-0001-7660-5055", r.scopusId = "55221928300", r.researcherId = "R-1968-2017"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "8609092415"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Illia", r.lastName = "Zymak", r.identificationNumber = "8609092415", r.orcid = "0000-0003-0224-7616", r.scopusId = "54401884200", r.researcherId = "P-6900-2018"
WITH r
MATCH (c:Country {code: "CZ"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);

MERGE (r:Researcher {identificationNumber: "911118"})
ON CREATE SET r.uid = apoc.create.uuid(), r.firstName = "Anna", r.lastName = "Zymaková", r.identificationNumber = "911118", r.orcid = "0000-0003-0175-638X"
WITH r
MATCH (c:Country {code: "UA"})
MERGE (r)-[:HAS_CITIZENSHIP]->(c);