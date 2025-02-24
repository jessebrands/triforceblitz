package com.triforceblitz.triforceblitz.randomizer;

/**
 * Domain implementation of the generator service.
 */
public class GeneratorServiceImpl implements GeneratorService {
    private final RandomizerRepository randomizerRepository;
    private final GeneratorFactory generatorFactory;

    public GeneratorServiceImpl(RandomizerRepository randomizerRepository,
                                GeneratorFactory generatorFactory) {
        this.randomizerRepository = randomizerRepository;
        this.generatorFactory = generatorFactory;
    }

    @Override
    public GeneratorOutput<Patch> generatePatch(RandomizerVersion version, String presetName, String seed) {
        var randomizer = randomizerRepository.load(version)
                .orElseThrow(() -> new RuntimeException("randomizer " + version + " not found"));

        var preset = randomizer.getPreset(presetName)
                .orElseThrow(() -> new RuntimeException("randomizer does not have preset " + presetName));

        return generatorFactory
                .create(randomizer, preset)
                .generatePatch(seed);
    }
}
