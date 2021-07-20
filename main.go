package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

// Baseado na implementação presente na página:
// https://en.wikipedia.org/wiki/Xorshift
func Xorshift(nBits int) *big.Int {
	retorno := big.NewInt(0)
	// Gerando seed usando o relógio
	x := uint32(time.Now().UnixNano())
	for retorno.BitLen() < nBits {
		// Cada iteração gera um aleatório de 32 bits e o concatena ao retorno
		x ^= x << 13
		x ^= x >> 17
		x ^= x << 5
		retorno.Lsh(retorno, uint(32))
		retorno.Or(retorno, big.NewInt(int64(x)))
	}
	// Removendo possíveis bits em excesso
	retorno.Rsh(retorno, uint(retorno.BitLen()-nBits))
	return retorno.Abs(retorno)
}

// Baseado na implementação presente na página:
// https://en.wikipedia.org/wiki/Linear_congruential_generator
func LCG(nBits int) *big.Int {
	retorno := big.NewInt(0)
	// Parametros escolhidos de acordo com Borland C/C++
	// Específicos para gerar números de 32 bits
	xN := big.NewInt(time.Now().UnixNano())
	xN.Abs(xN)
	a := big.NewInt(22695477)
	c := big.NewInt(1)
	m := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(32), nil)
	for retorno.BitLen() < nBits {
		// Cada iteração gera um aleatório de 32 bits e o concatena ao retorno
		xN.Mul(xN, a)
		xN.Add(xN, c)
		xN.Mod(xN, m)
		retorno.Lsh(retorno, uint(xN.BitLen()))
		retorno.Or(retorno, xN)
	}
	// Removendo possíveis bits em excesso
	retorno.Rsh(retorno, uint(retorno.BitLen()-nBits))
	return retorno
}

//Seja n um número primo
//Baseado no algoritmo presente na página:
//https://pt.wikipedia.org/wiki/Teste_de_primalidade_de_Miller-Rabin
//e sua implementação em Python presente em:
//https://gist.github.com/Ayrx/5884790
func MillerRabin(n *big.Int) bool {
	// Chance de falso positivo é 1/4^nTests
	nTests := 20
	// 2 e 3 são primos
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	// Se divisivel por 2, não é primo
	if big.NewInt(0).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	r := uint64(0)
	d := big.NewInt(0).Sub(n, big.NewInt(1))
	// Encontrando maior potencia de 2 que divide (d = n-1)
	for big.NewInt(0).Mod(d, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		r += uint64(1)
		d.Rsh(d, 1)
	}
	nMinus1 := big.NewInt(0).Sub(n, big.NewInt(1))

	for i := 0; i < nTests; i++ {
		// Usando um shift aleatório em n como número aleatório
		// a é um inteiro aleatório tal que 1 < a < n/2
		a := big.NewInt(0).Rsh(Xorshift(n.BitLen()), uint(rand.Intn(n.BitLen()-1)))
		// x = (a^d) mod n
		x := big.NewInt(0).Exp(a, d, n)
		// Caso x congruente a 1
		if (x.Cmp(big.NewInt(1)) == 0) || x.Cmp(nMinus1) == 0 {
			continue
		}
		possiblePrime := false
		// tentando encontrar possível inteiro tal que x = a^((2^r)) congruente a -1 mod n
		for j := uint64(0); j < (r - 1); j++ {
			// x = a^((2^r))
			x.Exp(x, big.NewInt(2), n)
			// Caso x seja congruente a -1 mod n, n é um possível primo
			if x.Cmp(nMinus1) == 0 {
				possiblePrime = true
				break
			}
		}
		if !possiblePrime {
			return false
		}
	}
	return true
}

// Checa se n é primo
// Baseado na aula/implementação presente em:
// https://pt.khanacademy.org/computing/computer-science/cryptography/random-algorithms-probability/pi/level-10-fermat-primality-test
// e na explicação fornecida em: https://pt.wikipedia.org/wiki/Teste_de_primalidade_de_Fermat
func Fermat(n *big.Int) bool {
	// Number of tests
	nTests := 20
	if n.Cmp(big.NewInt(1)) == 0 || n.Cmp(big.NewInt(4)) == 0 {
		return false
	} else if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	} else {
		for i := 0; i < nTests; i++ {
			nMinus1 := big.NewInt(0).Sub(n, big.NewInt(1))
			// Gerando número aleatório de 32 bits
			// a := Xorshift(32)
			// a é um inteiro aleatório tal que 1 < a < n/2
			// usando um número aleatório de shifts em cima de um aleatório de 1 bit a menos que n
			a := big.NewInt(0).Rsh(Xorshift(n.BitLen()), uint(rand.Intn(n.BitLen()-1)))
			if a.Exp(a, nMinus1, n).Cmp(big.NewInt(1)) != 0 {
				return false
			}
		}
	}
	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())
	BitsNeeded := []int{40, 56, 80, 128, 168, 224, 256, 512, 1024, 2048, 4096}

	fmt.Println("---- GERANDO NÚMEROS ALEATÓRIOS: LCG ----")
	for i := 0; i < len(BitsNeeded); i++ {
		start := time.Now()
		n := LCG(BitsNeeded[i])
		end := time.Now()
		time := end.Sub(start)
		fmt.Printf("Numero de %d bits gerado em %dµs:\n%v\n", n.BitLen(), time.Microseconds(), n)

	}

	fmt.Println("---- GERANDO NÚMEROS ALEATÓRIOS: XorShift ----")
	for i := 0; i < len(BitsNeeded); i++ {
		start := time.Now()
		n := Xorshift(BitsNeeded[i])
		end := time.Now()
		time := end.Sub(start)
		fmt.Printf("Numero de %d bits gerado em %dµs:\n%v\n", n.BitLen(), time.Microseconds(), n)
	}

	fmt.Println("---- GERANDO NÚMEROS PRIMOS: Fermat ----")
	for _, bits := range BitsNeeded {
		fmt.Printf("----  %d BITS ----\n", bits)
		var atempts int64 = 1
		a := big.NewInt(time.Now().UnixNano())
		a.Abs(a)
		startPrime := time.Now()
		for !Fermat(a) {
			a = Xorshift(bits)
			atempts++
		}
		endPrime := time.Now()
		timePrime := endPrime.Sub(startPrime)
		fmt.Printf("Gerado um primo com %d bits em %vµs após %d tentativas:\n%v\n", a.BitLen(), timePrime.Microseconds(), atempts, a)
	}
	fmt.Println("---- GERANDO NÚMEROS PRIMOS: Miller-Rabin ----")
	for _, bits := range BitsNeeded {
		fmt.Printf("----  %d BITS ----\n", bits)
		var atempts int64 = 1
		a := big.NewInt(time.Now().UnixNano())
		a.Abs(a)
		startPrime := time.Now()
		for !MillerRabin(a) {
			a = Xorshift(bits)
			atempts++
		}
		endPrime := time.Now()
		timePrime := endPrime.Sub(startPrime)
		fmt.Printf("Gerado um primo com %d bits em %vµs após %d tentativas:\n%v\n", a.BitLen(), timePrime.Microseconds(), atempts, a)
	}
}
