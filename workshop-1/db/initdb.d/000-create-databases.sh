#!/bin/sh
set -eu

. /config/databases.conf

# Проверка корректности идентификатора
validate_identifier() {
	case $1 in
		*[!a-zA-Z0-9_]*|"")
			echo "$(basename "$0"): Invalid identifier '$1'" >&2
			return 1
			;;
		*) 
			return 0 
			;;
	esac
}

var_is_empty() {
	eval test -z \"\${${1}-}\"
}

get_db_var() {
	local var_name="db_${1}_${2}"
	if var_is_empty "$var_name"; then
		echo "$(basename "$0"): $var_name not set" >&2
		return 1
	fi
	eval $2=\${$var_name}
}

# Экранирование апострофов для SQL
sql_escape() {
	echo "$1" | sed "s/'/''/g"
}

main() {
	# Проверка обязательной переменной
	: ${db_list_to_create:?Environment variable db_list_to_create not set}
	
	for db in $db_list_to_create; do
		# Валидация имени шаблона
		validate_identifier "$db" || continue

		# Чтение переменных
		get_db_var $db name     || continue
		get_db_var $db user     || continue
		get_db_var $db password || continue

		validate_identifier "$name" || continue
		validate_identifier "$user" || continue
		escaped_password=$(sql_escape "$password")

		echo "Processing: $name" >&2
		
		# Безопасный SQL через HEREDOC
		psql -U postgres -v ON_ERROR_STOP=1 <<-SQL
			SELECT 'CREATE DATABASE "${name}"' 
			WHERE NOT EXISTS (
				SELECT FROM pg_database 
				WHERE datname = '${name}'
			)\\gexec

			\\connect ${name}
			
			DO \$\$
			BEGIN
				IF NOT EXISTS (
					SELECT FROM pg_roles 
					WHERE rolname = '${user}'
				) THEN
					EXECUTE format(
						'CREATE USER %I WITH PASSWORD %L',
						'${user}',
						'${escaped_password}'
					);

					-- Отзываем публичные права для безопасности
					REVOKE ALL ON SCHEMA public FROM PUBLIC;

					-- Сделать пользователя владельцем базы
					ALTER DATABASE "${name}" OWNER TO "${user}";

					-- Сделать пользователя владельцем схемы public
					ALTER SCHEMA public OWNER TO "${user}";
				END IF;
			END
			\$\$;
		SQL
	done
}

main "$@"
