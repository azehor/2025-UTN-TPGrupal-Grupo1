-- Seed de juegos (basado en frontend/src/mock/juegosMock.ts)
-- Usa imágenes locales ya copiadas en: api/static/
-- Asegúrate de que los archivos existan con los nombres usados abajo.
-- Si la tabla se llama distinto, ajusta: softwares

BEGIN;

-- (Opcional) limpiar solo juegos previos
-- DELETE FROM softwares WHERE tipo = 'juego';

-- Asegurar extensión para gen_random_uuid (si no existe)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Inserciones:
-- Nota sobre image_url:
-- Ahora se sirven archivos directamente desde /static.
-- En DB guardamos solo el nombre de archivo (sin rutas) y el backend arma /static/<filename>.

INSERT INTO softwares
  (id, nombre, tipo, empresa, image_url, orden_grafica, orden_procesador, orden_ram)
VALUES
  (gen_random_uuid(),'Cyberpunk 2077','juego','CD Projekt Red','images/juegos/Cyberpunk_2077_box_art.jpg',9,8,8),
  -- El archivo en disco es "read.avif" para Red Dead Redemption 2
  (gen_random_uuid(),'Red Dead Redemption 2','juego','Rockstar Games','images/juegos/read.avif',10,9,8),
  (gen_random_uuid(),'Elden Ring','juego','FromSoftware','images/juegos/eldenring.avif',8,8,7),
  (gen_random_uuid(),'Half-Life: Alyx','juego','Valve','images/juegos/alix.jpg',10,9,9),
  (gen_random_uuid(),'Minecraft','juego','Mojang','images/juegos/Minecraft_2024_cover_art.png',3,4,4),
  (gen_random_uuid(),'Assassin''s Creed Valhalla','juego','Ubisoft','images/juegos/Assassins_Creed_Valhalla_capa.png',9,8,8),
  (gen_random_uuid(),'The Elder Scrolls V: Skyrim','juego','Bethesda','images/juegos/The_Elder_Scrolls_V_Skyrim_cover.png',8,7,7),
  (gen_random_uuid(),'The Witcher 3: Wild Hunt','juego','CD Projekt Red','images/juegos/Witcher_3_cover_art.jpg',9,8,8),
  (gen_random_uuid(),'Grand Theft Auto V','juego','Rockstar Games','images/juegos/Grand_Theft_Auto_V.png',9,8,8),
  (gen_random_uuid(),'Portal 2','juego','Valve','images/juegos/Portal2cover.jpg',7,6,6),
  (gen_random_uuid(),'Fortnite','juego','Epic Games','images/juegos/Fortnite_Save_The_World.jpg',6,5,5),
  (gen_random_uuid(),'The Last of Us Part II','juego','Naughty Dog','images/juegos/TLOU_P2_Box_Art_2.png',9,8,8),
  (gen_random_uuid(),'Resident Evil Village','juego','Capcom','images/juegos/Resident_Evil_Village.png',9,8,7),
  (gen_random_uuid(),'Final Fantasy VII Remake','juego','Square Enix','images/juegos/FFVIIRemake.png',8,7,7),
  (gen_random_uuid(),'Gwent','juego','CD Projekt Red','images/juegos/Gwent_cover_art.png',5,4,4),
  (gen_random_uuid(),'Dark Souls III','juego','FromSoftware','images/juegos/Dark_souls_3_cover_art.jpg',8,7,7),
  (gen_random_uuid(),'Left 4 Dead 2','juego','Valve','images/juegos/Left4Dead2.jpg',7,6,6),
  (gen_random_uuid(),'Minecraft Dungeons','juego','Mojang','images/juegos/Minecraft_Dungeons_cover_art.jpg',6,5,5),
  (gen_random_uuid(),'Far Cry 6','juego','Ubisoft','images/juegos/Far_cry_6_cover.jpg',9,8,7),
  (gen_random_uuid(),'Destiny 2','juego','Bungie','images/juegos/Destiny_2_(artwork).jpg',8,7,7);

COMMIT;

-- Verificación rápida:
-- SELECT nombre, image_url FROM softwares WHERE tipo='juego' ORDER BY nombre;