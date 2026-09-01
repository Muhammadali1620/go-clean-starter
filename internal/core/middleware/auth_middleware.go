package middleware

// // RequireAuth extracts the JWT, fetches the user, and sets the full user object in the context.
// func RequireAuth(jwtService services.JwtService, userRepo repositories.UserRepository) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			authHeader := c.Request().Header.Get("Authorization")
// 			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 				return models.NewAppError(models.ErrorTypeUnauthorized, models.ErrCodeUnauthorized, "missing or malformed auth header")
// 			}

// 			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 			claims, err := jwtService.ValidateToken(tokenString)
// 			if err != nil {
// 				return models.NewAppError(models.ErrorTypeUnauthorized, models.ErrCodeUnauthorized, "invalid or expired token")
// 			}

// 			// Fetch user from DB
// 			user, err := userRepo.FindByID(c.Request().Context(), claims.UserID)
// 			if err != nil {
// 				var appErr *models.AppError
// 				if errors.As(err, &appErr) && appErr.Type == models.ErrorTypeNotFound {
// 					return models.NewAppError(models.ErrorTypeUnauthorized, models.ErrCodeUnauthorized, "user not found or deleted")
// 				}
// 				return models.NewAppErrorWrap(models.ErrorTypeInternal, models.ErrCodeInternalError, "failed to authenticate, database is busy", err)
// 			}

// 			if user == nil {
// 				return models.NewAppError(models.ErrorTypeUnauthorized, models.ErrCodeUnauthorized, "user not found or deleted")
// 			}

// 			// Global check for banned users
// 			if !user.IsActive {
// 				return models.NewAppError(models.ErrorTypeForbidden, models.ErrCodeAccountBlocked, "user account is blocked")
// 			}

// 			// Save user in context
// 			c.Set("user", user)

// 			return next(c)
// 		}
// 	}
// }
