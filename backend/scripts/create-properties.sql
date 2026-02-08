-- Create properties for REPS V2
-- Run this script against your database to create the two properties

-- Property 1: 3522 Wood Bridge Dr
INSERT INTO properties (name, address, property_type, default_mileage)
SELECT '3522 Wood Bridge Dr', '3522 Wood Bridge Dr', 'residential', 0
WHERE NOT EXISTS (
    SELECT 1 FROM properties WHERE address = '3522 Wood Bridge Dr'
);

-- Property 2: 3500 Saindon St
INSERT INTO properties (name, address, property_type, default_mileage)
SELECT '3500 Saindon St', '3500 Saindon St', 'residential', 0
WHERE NOT EXISTS (
    SELECT 1 FROM properties WHERE address = '3500 Saindon St'
);

-- Verify properties were created
SELECT id, name, address, property_type FROM properties WHERE address IN ('3522 Wood Bridge Dr', '3500 Saindon St');
