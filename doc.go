package main

/*
mm hh jj MMM JJJ tâche
mm représente les minutes (de 0 à 59)
hh représente l'heure (de 0 à 23)
jj représente le numéro du jour du mois (de 1 à 31)
MMM représente l'abréviation du nom du mois (jan, feb, ...) ou bien le numéro du mois (de 1 à 12)
JJJ représente l'abréviation du nom du jour ou bien le numéro du jour dans la semaine :
0 = Dimanche
1 = Lundi
2 = Mardi
...
6 = Samedi
7 = Dimanche (représenté deux fois pour les deux types de semaine)
Pour chaque valeur numérique (mm, hh, jj, MMM, JJJ) les notations possibles sont :

* : à chaque unité (0, 1, 2, 3, 4...)
5,8 : les unités 5 et 8
2-5 : les unités de 2 à 5 (2, 3, 4, 5)
*\/3 : toutes les 3 unités (0, 3, 6, 9...)
10-20/3 : toutes les 3 unités, entre la dixième et la vingtième (10, 13, 16, 19)


Raccourcis	Description	Équivalent
@reboot	Au démarrage	Aucun
@yearly	Tous les ans	0 0 1 1 *
@annually	Tous les ans	0 0 1 1 *
@monthly	Tous les mois	0 0 1 * *
@weekly	Toutes les semaines	0 0 * * 0
@daily	Tous les jours	0 0 * * *
@midnight	Toutes les nuits	0 0 * * *
@hourly	Toutes les heures	0 * * * *


0 0 13 * 5			each friday and 13 of each month at 00 h 00
30 23 * * * 		each day at 23 h 30
5 * * * * 			each hour past 5 minutes
30 23 1 * * 		each 1 of the month at 23 h 30
28 22 * * 1 		each monday at 22 h 28
22 11 13 * 5 		each friday and each 13 of the month at 11 h 22
12 10 2-5 * * 		from 2 to 5 of each month at 10 h 12
59 23 *\/2 * * 		each even day of each month at 23 h 59 :
0 22 * * 1-5 		from Monday to Friday at 22 h 00
*\/5 * * * * 		every 5 minutes
*/
