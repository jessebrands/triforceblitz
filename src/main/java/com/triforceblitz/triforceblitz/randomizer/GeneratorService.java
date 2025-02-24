package com.triforceblitz.triforceblitz.randomizer;

/**
 * Services for generating Ocarina of Time ROMs.
 */
public interface GeneratorService {
    /**
     * Generates a new Ocarina of Time Randomizer ROM patch file.
     *
     * @param version Version of the Randomizer.
     * @param preset  Name of the preset to use.
     * @param seed    Seed value to initialize the randomizer.
     * @return Output of the generator.
     */
    GeneratorOutput<Patch> generatePatch(RandomizerVersion version, String preset, String seed);
}
