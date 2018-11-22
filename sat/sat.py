from random import *
from threading import Thread, Lock
#import matplotlib.pyplot as plt
import time
import math
import re

TAM = 250
N = 250000
T0 = 1
TN = 0.9999
t = T0
listaRSItens = []
listaSAItens = []
listaRandomsearch = []
listaSA = []
num_thread = 10
mutex = Lock()

def gerarRandomList(x):
    lista = []
    for i in range(x):
        lista.append(randrange(2))
    return lista

def ler():
    lista = []
    F = open("uf250_01.cnf","r")
    for linha in F:
        l = re.search("(\-?[0-9]+)\s+(\-?[0-9]+)\s+(\-?[0-9]+) 0", linha)
        if l != None:
            lista.append((int(l.group(1)), int(l.group(2)), int(l.group(3))))
    return lista

def energia(inicial,listaCNF):
    contGeral = 0
    for cada in listaCNF:
        cont = 0
        for x in cada:
            if(x<0 and not inicial[abs(x) - 1]):
                contGeral += 1
		break;
            elif (x>0 and inicial[abs(x) - 1]):
                contGeral += 1
		break;


    return contGeral

def temperatura(i):
    A = float(N**-2) * math.log(T0/TN)
    return T0*math.e**(-A*i**2)

def vizinho(lista):
    x = randrange(len(lista))
    novaL = []
    for i in range(len(lista)):
        novaL.append(lista[i])
    if(lista[x]):
        novaL[x] = 0
    else:
        novaL[x] = 1
    return novaL

def randomsearch(s0,listaCNF):
    #global listaRSItens
    cont = 1
    candidato = s0
    melhorEnergia = energia(candidato,listaCNF)
    melhorCandidato = candidato
    #lista = [melhorEnergia]
    while(cont < N):
        candidato = vizinho(candidato)
        vizinhoE = energia(candidato,listaCNF)
        #lista.append(vizinhoE)
        if(melhorEnergia < vizinhoE):
            melhorCandidato = candidato
            melhorEnergia = vizinhoE
        cont += 1
    #if not (cont < N):
    #listaRSItens.append(lista)
    return melhorEnergia

def simuAnne(s0,listaCNF,id,chunk):
    #global listaSAItens
    candidato = s0

    cont = 1
    e = energia(s0,listaCNF)
    limit = (id*chunk)+chunk
    #lista = [e]
    while(True):
        proximo = vizinho(candidato)
        deltaE = energia(candidato,listaCNF) - energia(proximo,listaCNF)

        if(deltaE <= 0):
            candidato = proximo
        elif random() + float(randrange(0, 99)) < math.e ** (-deltaE/float(t)):
            candidato = proximo

        mutex.acquire()
        t = temperatura(cont)
        mutex.release()
        cont += 1
        #lista.append((energia(candidato,listaCNF)))
        if(t < TN or cont >= limit):
            #listaSAItens.append(lista)
            asfaf = energia(candidato,listaCNF)
            return

def media_dp(lista):
    media = 0.0
    dp = 0.0
    for cada in lista:
        media += float(cada)

    for i in len(media):
        dp += float((cada - media)**2)

    media = media/N
    dp = math.sqrt(dp/N)
    return (media,dp)

listaCNF = ler()



chunk = N/8
inicial = gerarRandomList(TAM)
for i in range(8):
    x = Thread(target = simuAnne, args = (inicial,listaCNF,i,chunk))
    x.start()

"""
for i in range(9):
    print "{} Iniciou!\n".format(i)
    inicial = gerarRandomList(TAM)
    print "INICIAL = {}".format(energia(inicial,listaCNF))
    print "Inical da {} Terminou!\n".format(i)
    listaRandomsearch.append(randomsearch(inicial,listaCNF))
    print "RS da {} Terminou!\n".format(i)
    print "##### RANDOM SEARCH #####"
    print listaRandomsearch
    print "##### MEDIA e DP  RANDOM SEARCH #####"
    print media_dp(listaRSItens[i])
    listaSA.append(simuAnne(inicial,listaCNF))
    print "{} Finalizou!\n".format(i)


    print "##### SIMULATED ANNELING #####"
    print listaSA


    print "##### MEDIA e DP  SIMULATED ANNELING #####"
    print media_dp(listaSAItens[i])

    #plotGraph(listaSAItens[0],'TESTE')
"""
