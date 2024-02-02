package ai

import (
	"fmt"
	"runtime"
	"testing"
)

var longText = `About a year ago I started a new job as a Staff Engineer with a larger FinTech. Plenty of upside, and I’m enjoying myself a lot, but the downside is that I’m obliged to work on a Mac. I kind of hate Macs. Not quite as bad as Windows, but damn I’m missing Linux in the workplace.

Note: If you plan on following in my footsteps, do note that I did all the work described here on my personal Linux machine over a weekend - it’s only the end result that I used on the actual work Mac

However, this did give me an excuse to indulge my taste for retro-computing by picking up an original Macintosh Extended Keyboard II - this beauty:
Macintosh Extended Keyboard II

Partly I chose it because it’s of the vintage that I occasionally used during my final year at university in the Mac lab - mostly while coding up a fairly hapless AI project in Macintosh Common Lisp (MCL).

Apple always did things a bit differently, so while the IBM PC derived machines had more-or-less standardised on the PS/2 based connectors, Macs had their own serial bus system for keyboards and mice. The PC stuff and the Mac stuff both used a four pin DIN connector just to ensure maximum confusion, and of course they were completely incompatible.

The nice thing about the Apple system was that you could daisy chain it - and the keyboards generally had a connector on both sides, making a rare concession to lefties who could thus connect the keyboard on the right and chain the mouse off of it on the left side of the keyboard if they chose.

Incidentally, if you were in a whimsical student mood, you could also chain all of the keyboards on a row of machines end to end resulting in five or six machines completely unconnected from their keyboard and the one at the end taking any and all input from the neighbours. Ok, it was kind of a dick move, but in a lab full of IT students it wasn’t too mean a trick and nobody was stymied for long.

The old connection system was called “Apple Desktop Bus” (ADB). These days Macs are firmly USB oriented for their minor peripherals, so getting an old-skool keyboard to talk to a spiffy modern M1 Mac requires a suitable adapter. There are a handful of commercially available adapters, but it’s more fun to build something and it’s the kind of project that’s actually within reach of even my meagre electronics skills.

The QMK Project is the basis for the device that I therefore built - it’s better known as a basis for building custom keyboards, but it also contains some code for building adapters. The wiring required is well documented in the README for the source folder for the adapter logic: https://github.com/qmk/qmk_firmware/tree/master/keyboards/converter/adb_usb

The ASCII art wiring diagram stolen from that is:

I think I ended up using a 4.7KΩ resistor - but I can’t remember why. Probably I saw it in some online instructions somewhere, but I’ve lost it now - my apologies to whichever reference I’m thus failing to cite here.

I picked up a cheap clone of a Pro Micro board on Amazon. The original and the clone are both based around the atmega32u4 and the pinouts are identical. The circuit consists of the board, a single resistor, and the DIN connector, so there’s really very little to do.

The one “gotcha” worth mentioning is that the PD0 pin mentioned in the QMK documentation corresponds to the D3 pin on the processor board, not the D0 pin as you might (I did) assume initially.

Burning the firmware was actually (marginally) harder than soldering it all together! However you should ignore the instructions on the README for the adapter logic - that’s a bit outdated and seems to be left over from the original TMK origin of the code. Instead, follow the instructions outlined in the main QMK documentation:

After cloning the QMK repository locally, the first steps were actually to install a qmk python command independently of that:

Note in the above that the .local/bin directory is not, on Ubuntu, in your path by default due to a long-standing Ubuntu bug. Running the qmk setup command will prompt you to install all the pre-requisites necessary for building the actual adapter logic.

With that complete you’re good to go to build the adapter-specific code from the root of the checked out QMK repository:

The default here is the keyboard map - you can alternatively build and compile in any weird mapping you choose.

The result of the build should be a file converter_adb_usb_rev1_default.hex and the next task is to get that loaded into the microcontroller.

Happily this requires no special steps and everything is auto-detected with the following command:

With the adapter logic flashed onto the microcontroller and everything wired up, you should see the adapter show up in the USB device list as QMK ADB to USB Keyboard Converter:

A minor gotcha to mention - I had a lot of initial trouble flashing the image, but I eventually realised that this was because the cable I was using was garbage. It came with some cheap crap audio gadget and was apparently barely good enough to use for charging and nothing else. All my issues went away when I swapped it for a good one!

I built my adapter on a bit of veroboard and I wanted to hide that away. A bit of measuring with calipers, some tinkering in Tinkercad, and sliiiightly more do-overs due to my horrible 3D design skills than I might have liked, and I had a design for a case to put it all in:
`

func TestSentenceEmbedding(t *testing.T) {
	sm, err := NewSentenceEmbeddingModel()
	if err != nil {
		t.FailNow()
	}

	defer sm.Close()

	embds, err := sm.Embeddings([]string{"I am a dog", "I am a cat"})

	if err != nil {
		t.FailNow()
	}

	fmt.Println("Cosine Close:", Cosine(embds[0].Vectors, embds[1].Vectors))

}

func TestSentenceEmbeddingNotClose(t *testing.T) {
	sm, err := NewSentenceEmbeddingModel()
	if err != nil {
		t.FailNow()
	}

	defer sm.Close()

	embds, err := sm.Embeddings([]string{"Vectors are represented as plain old floating point slices, there are no special data types", "If you get an error that the mount flag isn't supported, that indicates that you either didn't enable buildkit with the above variable,"})

	if err != nil {
		t.FailNow()
	}

	fmt.Println("Cosine Not close:", Cosine(embds[0].Vectors, embds[1].Vectors))

}

func TestSentenceEmbeddingLong(t *testing.T) {
	sm, err := NewSentenceEmbeddingModel()
	if err != nil {
		t.FailNow()
	}

	defer sm.Close()

	embds, err := sm.Embeddings([]string{longText})

	if err != nil {
		t.FailNow()
	}

	fmt.Println(embds)

}

func BenchmarkSentenceEmbeddingLong(b *testing.B) {
	sm, err := NewSentenceEmbeddingModel()
	if err != nil {
		b.FailNow()
	}

	defer sm.Close()

	b.ResetTimer()

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		_, err := sm.Embeddings([]string{longText})

		if err != nil {
			b.FailNow()
		}
	}
	runtime.GC()
	runtime.ReadMemStats(&m2)
	b.Log("alloc:", m2.Alloc-m1.Alloc)
	b.Log("mallocs:", m2.Mallocs-m1.Mallocs)

}

func TestKeywordExtraction(t *testing.T) {
	sm, err := NewKeywordExtractionModel()
	if err != nil {
		t.FailNow()
	}

	defer sm.Close()

	embds, err := sm.Extract([]string{longText})

	if err != nil {
		t.FailNow()
	}

	fmt.Println(embds)

}

func BenchmarkKeywordExtraction(b *testing.B) {
	sm, err := NewKeywordExtractionModel()
	if err != nil {
		b.FailNow()
	}

	defer sm.Close()

	b.ResetTimer()

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		_, err := sm.Extract([]string{longText})

		if err != nil {
			b.FailNow()
		}
	}
	runtime.GC()
	runtime.ReadMemStats(&m2)
	b.Log("alloc:", m2.Alloc-m1.Alloc)
	b.Log("mallocs:", m2.Mallocs-m1.Mallocs)

}

func TestSingleShot(t *testing.T) {
	sm, err := NewZeroShotModel()
	if err != nil {
		t.FailNow()
	}

	defer sm.Close()

	embds, err := sm.Predict([]string{longText}, []string{"technology", "science", "linux", "life"})

	if err != nil {
		t.FailNow()
	}

	fmt.Println(embds)

}

func BenchmarkSingleShot(b *testing.B) {
	sm, err := NewZeroShotModel()
	if err != nil {
		b.FailNow()
	}

	defer sm.Close()

	b.ResetTimer()

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		_, err := sm.Predict([]string{longText}, []string{"technology", "science", "linux", "life"})

		if err != nil {
			b.FailNow()
		}
	}
	runtime.GC()
	runtime.ReadMemStats(&m2)
	b.Log("alloc:", m2.Alloc-m1.Alloc)
	b.Log("mallocs:", m2.Mallocs-m1.Mallocs)

}
