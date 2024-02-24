-- Create table for books
CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  author VARCHAR(255) NOT NULL,
  genre VARCHAR(100)
);

-- Generate 1000 books with random titles, authors, and genres
INSERT INTO books (title, author, genre)
SELECT 
  concat('Book ', generate_series(1, 1000)),
  concat('Author ', generate_series(1, 1000)),
  (array['Fiction', 'Non-Fiction', 'Biography', 'Science Fiction', 'Mystery'])[floor(random() * 5 + 1)]
FROM generate_series(1, 1000);

-- Create table for reviews
CREATE TABLE reviews (
  id SERIAL PRIMARY KEY,
  book_id INTEGER REFERENCES books(id) NOT NULL,
  rating INTEGER CHECK (rating BETWEEN 1 AND 5),
  comment VARCHAR(500)
);

-- Generate 5 reviews for each book with random ratings and comments
INSERT INTO reviews (book_id, rating, comment)
SELECT 
  b.id,
  floor(random() * 5) + 1,
  concat('This is a review for book ', b.title, '. It was ', 
  case 
    when floor(random() * 2) = 0 then 'good'
    else 'bad'
  end, '.')
FROM books b
CROSS JOIN generate_series(1, 5);
