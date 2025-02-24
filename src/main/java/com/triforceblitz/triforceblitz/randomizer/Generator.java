package com.triforceblitz.triforceblitz.randomizer;

/**
 * A generator is a configured Randomizer that can generate patched ROMs and
 * patch files.
 *
 * @see Randomizer
 * @see GeneratorFactory
 */
public interface Generator {
    GeneratorOutput<Patch> generatePatch(String seed);
}
