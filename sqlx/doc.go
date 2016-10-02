package sqlx

// Builder interface:
//
// Select(p ...interface{}) *Builder
// SelectRaw(exp string, bindings ...interface{}) *Builder
// Update(sqlx.Data) *Builder
// Insert(...sqlx.Data) *Builder
// Delete() *Builder
// From(table interface{}) *Builder
// Join(table, column1, operator, column2 string) *Builder
// LeftJoin(table, column1, operator, column2 string) *Builder
// Where(column string, operator string, value interface{}) *Builder
// WhereRaw(exp string, bindings ...interface{}) *Builder
// WhereGroup(callback func(*sqlx.Builder)) *Builder
// WhereIn(column string, values interface{}) *Builder
// WhereNotIn(column string, values interface{}) *Builder
// WhereBetween(column string, min, max interface{}) *Builder
// WhereNotBetween(column string, min, max interface{}) *Builder
// WhereNull(column string) *Builder
// WhereNotNull(column string) *Builder
// OrWhere(column string, operator string, value interface{}) *Builder
// OrWhereRaw(exp string, bindings ...interface{}) *Builder
// OrWhereGroup(callback func(*sqlx.Builder)) *Builder
// OrWhereIn(column string, values interface{}) *Builder
// OrWhereNotIn(column string, values interface{}) *Builder
// OrWhereBetween(column string, min, max interface{}) *Builder
// OrWhereNotBetween(column string, min, max interface{}) *Builder
// OrWhereNull(column string) *Builder
// OrWhereNotNull(column string) *Builder
// GroupBy(p ...interface{}) *Builder
// GroupByRaw(exp string, bindings ...interface{}) *Builder
// Having(column string, operator string, value interface{}) *Builder
// HavingRaw(exp string, bindings ...interface{}) *Builder
// HavingGroup(callback func(*sqlx.Builder)) *Builder
// OrHaving(column string, operator string, value interface{}) *Builder
// OrHavingRaw(exp string) *Builder
// OrHavingGroup(callback func(*sqlx.Builder)) *Builder
// OrderBy(column string, direction string) *Builder
// Limit(number interface{}) *Builder
// Offset(number interface{}) *Builder
// Count(column string, alias string) *Builder
// Sum(column string, alias string) *Builder
// Min(column string, alias string) *Builder
// Max(column string, alias string) *Builder
// Avg(column string, alias string) *Builder
//
// Sql() string
// Data() []interface{}
//
// Examples:
//
// query.Sql(); // строка запроса
// query.Data(); // срез с данными
//
// sqlx.Driver("mysql")
//
// SELECT * FROM `users`;
// query := sqlx.Table("users")
//
// SELECT `id` FROM `users`;
// query := sqlx.Table("users").Select("id")
//
// SELECT `id` FROM (SELECT `group_id`, MAX(`created_at`) as `lastdate` FROM `users` GROUP BY `group_id`) as `users` ORDER BY `lastdate` DESC;
// query := Table(func(builder *Builder) {
//     builder.Select("group_id").From("users").GroupBy("group_id").Max("created_at", "lastdate")
// }).OrderBy("lastdate", "DESC").Sql()
//
// SELECT * FROM (SELECT * FROM `users` WHERE `id` = ?) as users
// subque := Table("users").Where("id", "=", 1)
// query := Table(Raw("( "+subque.Sql()+" ) as users", subque.Data()...))
//
// SELECT * FROM (SELECT * FROM users WHERE id IN (?, ?, ?)) as users WHERE `age` > ?
// query := Table(Raw("(SELECT * FROM users WHERE id IN (?, ?, ?)) as users", 1, 2, 3)).Where("age", ">", 21)
//
// SELECT CONCAT('#', id) FROM `users`;
// query := sqlx.Table("users").Select(sqlx.Raw("CONCAT('#', id)"))
//
// SELECT CONCAT('#', id) FROM `users`;
// query := sqlx.Table("users").SelectRaw("CONCAT('#', id)")
//
// SELECT * FROM `users` WHERE `id` = ?;
// query := sqlx.Table("users").Where("id", "=", 1).Select("*")
//
// SELECT * FROM `users` WHERE `age` = ? OR `age` = ?;
// query := sqlx.Table("users").Where("age", "=", 18).OrWhere("age", "=", 21).Select("*")
//
// SELECT * FROM `users` WHERE `age` = ? AND (`city` = ? OR `city` = ?);
// query := sqlx.Table("users").
//	  Where("age", "=", 18).
//	  WhereGroup(func(builder *sqlx.Builder) {
//		  builder.Where("city", "=", 78).OrWhere("city", "=", 98)
//	  }).Select("*")
//
// SELECT * FROM `users` WHERE id = 1;
// query := sqlx.Table("users").WhereRaw("id = 1").Select("*")
//
// SELECT * FROM `users` WHERE id = 1;
// query := sqlx.Table("users").Where("id", "=", sqlx.Raw("1")).Select("*")
//
// SELECT * FROM `users` WHERE `id` IN ( ?, ?, ?, ?, ?, ? );
// query := sqlx.Table("users").WhereIn("id", sqlx.List{1, 2, 3, 4, 5, 6}).Select("*")
//
// SELECT * FROM `users` WHERE `id` IN (SELECT `users_id` FROM `order` WHERE `type` = ? );
// query := sqlx.Table("users").WhereIn("id", func(builder *sqlx.Builder) {
//		  builder.Select("users_id").From("order").Where("type", "=", 1)
//	  }).Select("*")
//
// SELECT * FROM `users` WHERE `id` BETWEEN ? AND ?;
// query := sqlx.Table("users").WhereBetween("id", 1, 10).Select("*")
//
// SELECT * FROM `users` WHERE `address` IS NULL;
// query := sqlx.Table("users").WhereNull("address").Select("*")
//
// SELECT * FROM `users` as `us` INNER JOIN `info` as `inf` ON (`inf`.`id` = `us`.`users_id`);
// query := sqlx.Table("users as us").Join("info as inf", "inf.id", "=", "us.users_id").Select("*")
//
// SELECT * FROM `users` as `us` LEFT JOIN `info` as `inf` ON (`us`.`id` = `inf`.`users_id`);
// query := sqlx.Table("users as us").LeftJoin("info as inf", "us.id", "=", "inf.users_id").Select("*")
//
// SELECT * FROM `users` LIMIT ? OFFSET ?;
// query := sqlx.Table("users").Limit(10).Offset(50).Select("*")
//
// SELECT * FROM `users` ORDER BY `name` DESC;
// query := sqlx.Table("users").OrderBy("name", "DESC").Select("*")
//
// SELECT "сity", "created_at" FROM `users` GROUP BY `сity`, `created_at`;
// query := sqlx.Table("users").GroupBy("сity", "created_at").Select("сity", "created_at")
//
// SELECT COUNT(*) as count FROM `users` GROUP BY `sity` HAVING `count` > ?;
// query := sqlx.Table("users").GroupBy("sity").Having("count", ">", 5).Select(sqlx.Raw("COUNT(*) as count"))
//
// SELECT COUNT(`id`) as `count` FROM `users`;
// query := sqlx.Table("users").Count("id", "count")
//
// SELECT SUM(`age`) as `sum_age` FROM users;
// query := sqlx.Table("users").Sum("age", sum_age)
//
// SELECT SUM(sex='man') as `sum_man`, SUM(sex='woman') as `sum_woman` FROM users;
// query := sqlx.Table("users").Sum(sqlx.Raw("sex='man'"), sum_man).Sum(sqlx.Raw("sex='woman'"), sum_woman)
//
// UPDATE `users` SET `name` = ? WHERE `id` = ?
// query := sqlx.Table("users").Where("id", "=", 1).Update(sqlx.Data{"name": "Jack"})
//
// INSERT INTO `users` (`id`, `name`) VALUES (?, ?)
// query := sqlx.Table("users").Insert(sqlx.Data{"id": 1, "name": "Jack"})
//
// INSERT INTO `users` (`id`, `name`) VALUES (?, ?), (?, ?)
// query := sqlx.Table("users").Insert(sqlx.Data{"id": 1, "name": "Jack"}, sqlx.Data{"id": 1, "name": "Make"})
//
// INSERT INTO `users` (`id`, `name`) VALUES (?, ?), (?, ?)
// query := sqlx.Table("users").Insert(sqlx.Data{"id": 1, "name": "Jack"}).Insert(sqlx.Data{"id": 1, "name": "Make"})
//
// DELETE FROM `users` WHERE `id` = ?
// query := sqlx.Table("users").Where("id", "=", 1).Delete()
