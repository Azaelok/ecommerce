Comando para iniciar postgres
-----
pg_ctl -D "C:\Program Files\PostgreSQL\16\data" start


-- Instalación de SODA
go install github.com/gobuffalo/pop/v6/soda@latest



soda.exe generate fizz CreateRoomRestrictionsTable
 *...UP
 *...Down
soda.exe migrate
soda.exe migrate down


soda generate fizz AddNotNullToReservationIDForRestrictions

    Esto es para mofidicar una columna de la tabla
    -----------------------------------------------
    change_column("room_restrictions", "restriction_id","integer",{"null": true})


soda reset - No nos funciono


go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool@v5.7.1


// Generamos el formato como SQL
soda generate sql SeedRoomsTable
    * Inserción de las habitaciones disponibles


// Se genera el formato SQL para 
soda generate sql SeedRestrictionsTable
    *Inserción en la tabla de restricciones


