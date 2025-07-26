package auth0_test

import (
	"context"
	"testing"

	middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/auth0"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_CustomAuth0Claims_Validate(t *testing.T) {
	Convey("When validating custom Auth0 claims", t, func() {
		claims := middleware.CustomAuth0Claims{
			Scope: "read:all write:users admin:system",
		}

		Convey("Should always return nil error", func() {
			err := claims.Validate(context.Background())
			So(err, ShouldBeNil)
		})

		Convey("Should validate with different contexts", func() {
			err := claims.Validate(context.TODO())
			So(err, ShouldBeNil)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			err = claims.Validate(ctx)
			So(err, ShouldBeNil)
		})
	})
}

func Test_CustomAuth0Claims_HasAllScopes(t *testing.T) {
	Convey("When checking if claims have all required scopes", t, func() {
		Convey("With valid scopes", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "read:all write:users admin:system",
			}

			Convey("Should return true for single existing scope", func() {
				result := claims.HasAllScopes([]string{"read:all"})
				So(result, ShouldBeTrue)
			})

			Convey("Should return true for multiple existing scopes", func() {
				result := claims.HasAllScopes([]string{"read:all", "write:users"})
				So(result, ShouldBeTrue)
			})

			Convey("Should return true for all scopes", func() {
				result := claims.HasAllScopes([]string{"read:all", "write:users", "admin:system"})
				So(result, ShouldBeTrue)
			})

			Convey("Should return false for non-existing scope", func() {
				result := claims.HasAllScopes([]string{"delete:everything"})
				So(result, ShouldBeFalse)
			})

			Convey("Should return false for mix of existing and non-existing scopes", func() {
				result := claims.HasAllScopes([]string{"read:all", "non:existent"})
				So(result, ShouldBeFalse)
			})

			Convey("Should return true for empty scope list", func() {
				result := claims.HasAllScopes([]string{})
				So(result, ShouldBeTrue)
			})
		})

		Convey("With empty scope", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "",
			}

			Convey("Should return false for any required scope", func() {
				result := claims.HasAllScopes([]string{"read:all"})
				So(result, ShouldBeFalse)
			})

			Convey("Should return true for empty scope list", func() {
				result := claims.HasAllScopes([]string{})
				So(result, ShouldBeTrue)
			})
		})

		Convey("With single scope", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "read:all",
			}

			Convey("Should return true for that scope", func() {
				result := claims.HasAllScopes([]string{"read:all"})
				So(result, ShouldBeTrue)
			})

			Convey("Should return false for different scope", func() {
				result := claims.HasAllScopes([]string{"write:users"})
				So(result, ShouldBeFalse)
			})
		})

		Convey("With scope containing extra spaces", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "read:all  write:users   admin:system",
			}

			Convey("Should handle extra spaces correctly", func() {
				result := claims.HasAllScopes([]string{"read:all", "write:users", "admin:system"})
				So(result, ShouldBeTrue)
			})
		})

		Convey("With scope has leading or trailing spaces", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: " read:all write:users ",
			}

			Convey("Should handle leading/trailing spaces correctly", func() {
				result := claims.HasAllScopes([]string{"read:all", "write:users"})
				So(result, ShouldBeTrue)
			})
		})

		Convey("With partial scope matches", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "read:all write:users",
			}

			Convey("Should not match partial scope names", func() {
				result := claims.HasAllScopes([]string{"read"})
				So(result, ShouldBeFalse)

				result = claims.HasAllScopes([]string{"write"})
				So(result, ShouldBeFalse)

				result = claims.HasAllScopes([]string{"read:a"})
				So(result, ShouldBeFalse)
			})
		})

		Convey("With complex scope scenarios", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "scope:1 scope:2 scope:1 duplicate:scope",
			}

			Convey("Should handle duplicate scopes correctly", func() {
				result := claims.HasAllScopes([]string{"scope:1"})
				So(result, ShouldBeTrue)

				result = claims.HasAllScopes([]string{"duplicate:scope"})
				So(result, ShouldBeTrue)
			})
		})

		Convey("With special characters in scopes", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "read:users-admin write:data.edit delete:all_items",
			}

			Convey("Should handle scopes with special characters", func() {
				result := claims.HasAllScopes([]string{"read:users-admin"})
				So(result, ShouldBeTrue)

				result = claims.HasAllScopes([]string{"write:data.edit"})
				So(result, ShouldBeTrue)

				result = claims.HasAllScopes([]string{"delete:all_items"})
				So(result, ShouldBeTrue)

				result = claims.HasAllScopes([]string{"read:users", "admin"})
				So(result, ShouldBeFalse)
			})
		})

		Convey("With single space between scopes", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "a b",
			}

			Convey("Should properly split on single spaces", func() {
				result := claims.HasAllScopes([]string{"a"})
				So(result, ShouldBeTrue)

				result = claims.HasAllScopes([]string{"b"})
				So(result, ShouldBeTrue)

				result = claims.HasAllScopes([]string{"a", "b"})
				So(result, ShouldBeTrue)
			})
		})

		Convey("With multiple consecutive spaces", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "scope1     scope2",
			}

			Convey("Should handle multiple spaces correctly", func() {
				result := claims.HasAllScopes([]string{"scope1", "scope2"})
				So(result, ShouldBeTrue)
			})
		})

		Convey("With empty string in required scopes", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "read:all write:users",
			}

			Convey("Should handle empty string in required scopes", func() {
				result := claims.HasAllScopes([]string{""})
				So(result, ShouldBeFalse)
			})
		})

		Convey("With tab and newline characters", func() {
			claims := middleware.CustomAuth0Claims{
				Scope: "scope1\tscope2\nscope3",
			}

			Convey("Should not split on tabs or newlines (only spaces)", func() {
				result := claims.HasAllScopes([]string{"scope1\tscope2\nscope3"})
				So(result, ShouldBeTrue)

				result = claims.HasAllScopes([]string{"scope1"})
				So(result, ShouldBeFalse)
			})
		})
	})
}
