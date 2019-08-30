.PHONY: clean-start start stop

clean-start:
	sudo rm -Rdf ./autonity-data*
	sudo rm -Rdf ./logs*.log
	touch logs1.log
	touch logs2.log
	touch logs3.log
	touch logs4.log
	touch logs5.log
	docker-compose up -d

start:
	docker-compose up -d

stop:
	docker-compose down
