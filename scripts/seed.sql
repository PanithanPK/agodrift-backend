-- Hotel booking schema seed

-- Drop old demo tables if they exist
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS hotels;
DROP TABLE IF EXISTS users;

-- Users table: stores admin and normal users
CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(50) NOT NULL DEFAULT 'user',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Hotels table: stores hotel inventory and data used for filters on the frontend
CREATE TABLE IF NOT EXISTS hotels (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  location VARCHAR(255) NOT NULL,          -- e.g. "Manhattan, New York"
  destination VARCHAR(255) NOT NULL,       -- normalized destination/city for searching
  rating DECIMAL(2,1) NOT NULL DEFAULT 0,  -- e.g. 4.8
  reviews INT NOT NULL DEFAULT 0,
  price_cents INT NOT NULL,                -- current price per night in cents
  original_price_cents INT NULL,           -- original price to show discount
  amenities TEXT,                          -- comma-separated list of amenities
  featured TINYINT(1) NOT NULL DEFAULT 0,  -- 1 = featured, 0 = normal
  max_adults INT NOT NULL DEFAULT 1,       -- capacity configuration
  max_children INT NOT NULL DEFAULT 0,
  rooms_total INT NOT NULL DEFAULT 1,      -- how many rooms this hotel has
  rooms_available INT NOT NULL DEFAULT 1,  -- rooms currently available for booking
  status VARCHAR(50) NOT NULL DEFAULT 'active', -- active / inactive / maintenance
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Bookings table: link between users and hotels with stay details
CREATE TABLE IF NOT EXISTS bookings (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  hotel_id INT NOT NULL,
  check_in DATE NOT NULL,
  check_out DATE NOT NULL,
  adults INT NOT NULL DEFAULT 1,
  children INT NOT NULL DEFAULT 0,
  rooms INT NOT NULL DEFAULT 1,
  total_price_cents INT NOT NULL,
  status VARCHAR(50) NOT NULL DEFAULT 'pending', -- confirmed / cancelled / pending
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_bookings_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_bookings_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id)
);

-- Seed users (including an admin)
INSERT INTO users (name, email, password, role) VALUES
('Admin User', 'admin@agodrift.dev', 'adminpass', 'admin'),
('Alice Traveler', 'alice@example.com', 'userpass', 'user');

-- Seed hotels based on frontend demo data
INSERT INTO hotels (
  name,
  description,
  location,
  destination,
  rating,
  reviews,
  price_cents,
  original_price_cents,
  amenities,
  featured,
  max_adults,
  max_children,
  rooms_total,
  rooms_available,
  status
) VALUES
(
  'The Azure Downtown Hotel',
  'Modern city hotel in the heart of Manhattan, perfect for business and leisure stays.',
  'Manhattan, New York',
  'New York',
  4.8,
  2847,
  18900,
  24900,
  'Free Wi-Fi,Pool,Breakfast',
  1,
  2,
  1,
  80,
  45,
  'active'
),
(
  'Oceanview Paradise Resort',
  'Beachfront resort with stunning ocean views and spa facilities.',
  'Maldives',
  'Maldives',
  4.9,
  1523,
  32000,
  NULL,
  'Ocean View,Spa,Restaurant',
  1,
  2,
  2,
  60,
  25,
  'active'
),
(
  'Alpine Mountain Lodge',
  'Cozy lodge with fireplace and ski access in the Swiss Alps.',
  'Zermatt, Switzerland',
  'Zermatt',
  4.7,
  892,
  27500,
  35000,
  'Ski Access,Fireplace,Parking',
  0,
  2,
  2,
  40,
  18,
  'active'
),
(
  'Tropical Island Villas',
  'Luxury overwater villas with private beach and butler service.',
  'Bora Bora, French Polynesia',
  'Bora Bora',
  4.9,
  1245,
  45000,
  NULL,
  'Private Beach,Pool,Butler Service',
  1,
  3,
  2,
  30,
  10,
  'active'
),
(
  'Villa Romantica',
  'Historic villa with sea views on the Amalfi Coast.',
  'Amalfi Coast, Italy',
  'Amalfi Coast',
  4.6,
  687,
  22500,
  28000,
  'Historic,Terrace,Breakfast',
  0,
  2,
  1,
  20,
  8,
  'active'
),
(
  'Metropolitan Grand Hotel',
  'Premium city hotel in central Tokyo with skyline views.',
  'Tokyo, Japan',
  'Tokyo',
  4.8,
  3421,
  21000,
  NULL,
  'City View,Gym,Concierge',
  0,
  2,
  1,
  120,
  70,
  'active'
);

-- Seed example bookings linking users and hotels
INSERT INTO bookings (
  user_id,
  hotel_id,
  check_in,
  check_out,
  adults,
  children,
  rooms,
  total_price_cents,
  status
) VALUES
(
  2,
  1,
  '2025-03-10',
  '2025-03-14',
  2,
  1,
  1,
  18900 * 4,
  'confirmed'
),
(
  2,
  4,
  '2025-06-01',
  '2025-06-05',
  2,
  0,
  1,
  45000 * 4,
  'pending'
);

-- Useful indexes for query performance

-- Users: lookup by email and role
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);

-- Hotels: searching and filtering
CREATE INDEX IF NOT EXISTS idx_hotels_destination ON hotels (destination);
CREATE INDEX IF NOT EXISTS idx_hotels_rating ON hotels (rating);
CREATE INDEX IF NOT EXISTS idx_hotels_price ON hotels (price_cents);
CREATE INDEX IF NOT EXISTS idx_hotels_featured_status ON hotels (featured, status);

-- Bookings: lookups and availability checks
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings (user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_hotel_id ON bookings (hotel_id);
CREATE INDEX IF NOT EXISTS idx_bookings_hotel_dates ON bookings (hotel_id, check_in, check_out);
