package storage

// func TestCreateAndGetUser_Success(t *testing.T) {
// 	ctx := context.TODO()

// 	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
// 		os.Getenv("DB_USER"),
// 		os.Getenv("DB_PASS"),
// 		os.Getenv("DB_HOST"),
// 		os.Getenv("DB_PORT"),
// 		os.Getenv("DB_NAME"))
// 	log.Println(dsn)
// 	conn, err := pgx.Connect(ctx, dsn)

// 	defer conn.Close(ctx)

// 	pg := NewStore(conn)

// 	newUser := models.User{
// 		Username:       "admin1",
// 		HashedPassword: "pas",
// 		Salt:           "salt",
// 	}
// 	actualUser, err := pg.CreateUser(ctx, newUser)
// 	require.NotEmpty(t, actualUser.ID)
// 	require.NoError(t, err)
// 	require.Equal(t, newUser.Username, actualUser.Username)

// }
