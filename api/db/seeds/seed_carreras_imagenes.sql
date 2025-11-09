-- Seed patch: asignar imagenes a carreras (almacenamos SOLO el filename)
-- Ejecutar en PostgreSQL contra la misma base usada por la API
-- Las imagenes deben existir en api/static y la API completará la URL en runtime (/static/<filename>)

BEGIN;

-- Desarrollo de Software -> carrera-ing-software.png
UPDATE carreras
SET image_url = 'images/carreras/carrera-ing-software.png'
WHERE id = 'f24f0164-0942-49d6-8962-43c8a5c5b620';

-- Diseño Gráfico -> carrera-disenio-grafico.png
UPDATE carreras
SET image_url = 'images/carreras/carrera-disenio-grafico.png'
WHERE id = 'f53086c5-68df-4cc5-a1cc-f4b55090e6e3';

-- Edición de Video -> carrera-edicion-video.png
UPDATE carreras
SET image_url = 'images/carreras/carrera-edicion-video.png'
WHERE id = 'c3645691-4afa-4933-a2da-e561c691735e';

-- Ingeniería Mecánica -> carrera-ing-mecanica.png
UPDATE carreras
SET image_url = 'images/carreras/carrera-ing-mecanica.png'
WHERE id = '0bdf4aab-0001-48cb-b886-f429dcb445d3';

-- Arquitectura -> carrera-arquitectura.png
UPDATE carreras
SET image_url = 'images/carreras/carrera-arquitectura.png'
WHERE id = 'f16b6a99-4814-4efc-b278-66ba3b2080ed';

-- Entretenimiento y Videojuegos -> carrera-entretenimiento-videojuegos.png
UPDATE carreras
SET image_url = 'images/carreras/carrera-entretenimiento-videojuegos.png'
WHERE id = '76b8d2b9-0bce-494e-82e7-6cbaaf34f4c9';

-- Derecho -> por ahora placeholder generico (reemplazar cuando se suba imagen especifica)
UPDATE carreras
SET image_url = 'images/carreras/carrera-derecho.png'
WHERE id = '1783b7f0-cc3f-40b1-8908-9f36926ea6b5';

COMMIT;
