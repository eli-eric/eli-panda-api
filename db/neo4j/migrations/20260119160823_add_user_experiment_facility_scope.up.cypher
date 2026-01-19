// Create UserExperiment constraint and index (if not exist)
CREATE CONSTRAINT userexperiment_uid_unique IF NOT EXISTS FOR (n:UserExperiment) REQUIRE n.uid IS UNIQUE;
CREATE INDEX userexperiment_name_index IF NOT EXISTS FOR (n:UserExperiment) ON (n.name);

// Add facility relationship to existing UserExperiment nodes (if any exist without it)
MATCH (ue:UserExperiment)
WHERE NOT (ue)-[:BELONGS_TO_FACILITY]->(:Facility)
MATCH (f:Facility {code: "B"})
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

// Seed data for ELI-Beamlines facility (75 entries)
MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-19"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-19 Developing Hadamard Time-resolved Crystallography at the Laser-driven Plasma X-ray source (HATRX@PXS)"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-22"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-22 Probing the dynamics of resonant ICD in helium nanodroplets doped with rubidium atoms and molecule"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-23"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-23 Coherent Diffractive Imaging and broadband Ptychography of biological specimen at MAC endstation utilizing ultra-short EUV pulses from laser-driven Higher-Harmonics Generation"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-24"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-24 Ultrafast excited-state dynamics of a dinuclear iridium complex"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-25"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-25 Delocalization of electronic excitation of the metal center in a blue copper protein"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-33"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-33 Towards Molecular Frame Photoelectron Angular Distribution Measurements of Polyatomic Molecules"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-39"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-39 High repetition rate EUV laser for applications"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-44"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-44 Optimization of TCT setup and method for timing and spatial characterization of novel 4D pixels, TI-LGADs, for future upgrades of CMS and ATLAS and for future colliders"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-49"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-49 Time-resolved ellipsometry for organic photovoltaics"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-50"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-50 Search for coherent phonon oscillations in the dielectric function dynamics of GaP using fs-time resolved ellipsometry"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-51"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-51 Resolving the Reaction Dynamics of Bestrhodopsin, a Unique Red-Shifted Microbial Rhodopsin and Light-Modulated Anion Channel"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-52"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-52 Stimulated Raman scattering as non-perturbing method to probe p53-FOXO4 interaction"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-56"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-56 Disulfide bond cleavage and recombination in a biomimetic dicopper-disulfide complex"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-57"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-57 Determination of ultrafast carrier dynamics in colloidal ZnSe quantum dots"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-59"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-59 Ultrafast measurements of the surface plasmon resonance dynamics of metal nanomaterials based on gold nanoparticles to investigate steady-state and time-resolved optical properties"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "UPM-64"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "UPM-64 Time-resolved spectroscopy of cascaded inter-atomic decays in alkali aggregates in He nanodroplets"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-032"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-032 Pump-Probe Time-Resolved Spectroscopic Ellipsometry on Ge and GeSn Alloys"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-076"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-076 Study of Gallium Sulfide Phase Change Material at the Femtosec scale"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-084"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-084 Time-resolved Ellipsometry on ultrawideband gap metal oxides"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-086"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-086 Formation and Recombination Processes of Self-Trapped Excitons in Halide Perovskites with Broadband Luminescence"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-087"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-087 High repetition rate EUV laser for applications"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-089"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-089 Exploring Macromolecular Dynamics with Plasma X-ray Sources: X-ray Powder Diffraction of Hemoglobin Microcrystal"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-091"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-091 Pump-Probe Time-Resolved Spectroscopic Ellipsometry on Ge and GeSn Alloys"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-092"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-092 Two-photon photophysics of tungsten arylisocyanide complexes"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-096"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-096 TIME RESOLVED ELLIPSOMETRY OF INGAAS QUANTUM WELLS"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-105"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-105 High-Repetition Rate X-Ray Source based on Laser Plasma Accelerator"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-110"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-110 High-Repetition-Rate regime for laser-triggered nuclear fusion reactions"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-113"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-113 Laser Wakefield Electron Acceleration using tailored plasma targets"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-117"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-117 Investigating the equation of state of boron nitride in extreme conditions by direct-laser drive"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-121"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-121 Photoneutron generation by Undepleted Direct Laser Acceleration"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-125"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-125 kHz all-optical Thomson scattering source for medical imaging"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-129"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-129 Time-resolved measurements on an efficient ICD induced by photoelectron impact excitation in pure and doped large He nanodroplets"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-132"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-132 Probing ultrafast structural dynamics in Indium- and Bismuth based metal halide perovskites"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-135"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-135 high-repetition rate, high flux positron sources for industrial applications"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-137"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-137 Demonstration of a compact, high-rep dose delivery system employing helical coil targets"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-140"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-140 Study of the Multiplication Dynamics in the Interpad Regions of the segmented LGADs within low and high intensity injection and proposal for upgrade of the existing TCT at ELI Beamlines to enable study on SiC based sensors"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-141"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-141 Stability of the contact discontinuity between an adiabatic and a radiative optically thin shock."
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-144"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-144 A new route to probing ultrafast relaxation dynamics and HHG of excited helium nanodroplets"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-146"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-146 High-Repetition Rate Broadband Betatron X-Ray Source based on Laser Plasma Accelerator"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-147"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-147 Optical and nonlinear characteristics of air and laser-created air plasma during four wave mixing in the infrared spectral range"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-148"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-148 Solid-liquid interfaces investigated by femtosecond time-resolved spectroscopic ellipsometry."
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-149"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-149 Resolving the structural dynamics of NeoR, a bistable red-shifted microbial rhodopsin employing unique electronic structures and isomeric states"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-154"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-154 In Vivo Preclinical Investigation Using Zebrafish Embryo model to Study the Biological Effectiveness of Laser-Driven Electron Beam"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-155"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-155 Site-specific Raman spectroscopy of photoactive proteins with novel genetically encoded alkynes"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-156"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-156 Optimizing high rep-rate radiography with machine learning"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-164"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-164 Characterizing the equation of state of polyethylene for mimicking planetary interiors and the synthesis of new carbon materials"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM-165"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM-165 HATRX@TREX"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-15"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-15 In-Air Laser-Driven Particle Induced X-ray Emission of artwork-like materials with advanced Double-Layer Targets at the ELIMAIA beamline"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-25"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-25 Investigation of the transient dielectric function of the transition metal Niobium using femtosecond time-resolved ellipsometry"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-26"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-26 Carrier Dynamics on GeSn Alloys"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-27"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-27 Experimental measurement and control of relativistic flying mirror acceleration"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-31"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-31 Laser-driven biomolecular condensate formation in protein solutions"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-32"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-32 Complete photoinduced folding of cytochrome c monitored by site-specific Raman spectroscopy"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-35"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-35 Revealing the Mechanism of Thermally Activated Delayed Fluorescence of Carbonyl-based Narrowband Emitters"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-36"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-36 Femtosecond spectroscopic ellipsometry (FSE) with two photon absorption (TPA) for indirect semiconductors"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-37"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-37 Transient absorption spectroscopic measurements on Fe(II)-polypyridyl complexes with up to 500 ns delay times"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-40"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-40 Redox-dependent excitation quenching and annihilation in chlorosomes of green photosynthetic bacteria"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-42"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-42 Time domain characterization and modeling of plasmonic gratings"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-44"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-44 Absolute dosimetry for the ELIMED beamline"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-45"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-45 Reaction dynamics of bistable ultraviolet and visible light absorbing crustacean and algal rhodopsins"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-48"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-48 Investigation of laser induced discoloration in dielectric coatings by time-resolved ellipsometry"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-50"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-50 Laser plasma interaction in the context of Inertial Confinement Fusion: towards the identification of the key parameters and their potential to mitigate LPI"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-51"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-51 Tryptophan radical cations in long-distance charge transport: A single-Trp azurin system"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-53"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-53 Tuneable and high-flux inverse Compton scattering source of gamma-rays for industrial applications"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-57"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-57 Femtosecond laser-based SPA-TCT and TPA-TCT study of NLGAD detector"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-58"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-58 Targeting radioresistance in 3D brain tumour models with laser accelerated protons from ELIMAIA"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-73"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-73 laSer drIven hadroN cAncer radioTherapy on bReast cAncer"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-75"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-75 Transient optical phenomena in organic semiconductors"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-77"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-77 Demonstration of an all-optical laser waveguide for LWFA in ELBA experimental platform"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-79"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-79 Intersystem Crossing and Electronic Structure of Novel Ni(II) Porphyrin-Nanographene Conjugates"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-87"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-87 Radiobiological Investigation on Laser-Driven Proton Beam"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-90"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-90 Using laser-accelerated protons and nuclear fusion reactions for widening the protontherapy therapeutic index"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM3-91"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM3-91 Exploring the Electric Field Thresholds and Conductivity Effects on Ghost multiplication in Interpad Region of TI-LGAD with Double Trenches and Study of Impact Ionization on dedicated LGAD sample"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM4-014"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM4-014 Transient Dielectric Function Of Squaraine Thin Films and ITO Electrodes"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM4-123"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM4-123 Uncovering new intermediates in the unique carotenoid-driven process using combined time-resolved spectral and crystallographic approaches."
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (ue:UserExperiment {code: "ELIUPM4-128"}) ON CREATE SET ue.uid = randomUUID(), ue.name = "ELIUPM4-128 Dechiphering the mechanism of bistability in opsins by structure and function"
MERGE (ue)-[:BELONGS_TO_FACILITY]->(f);

// Data migration: Create relationships from existing Publication.userExperiment strings
// Note: keeping the string property for backward compatibility
MATCH (p:Publication) WHERE p.userExperiment IS NOT NULL AND p.userExperiment <> ""
MATCH (ue:UserExperiment {name: p.userExperiment})
MERGE (p)-[:HAS_USER_EXPERIMENT]->(ue);
