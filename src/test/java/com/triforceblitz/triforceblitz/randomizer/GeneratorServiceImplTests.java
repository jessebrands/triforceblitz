package com.triforceblitz.triforceblitz.randomizer;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Map;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GeneratorServiceImplTests {
    @Mock
    private RandomizerRepository randomizerRepository;

    @Mock
    private GeneratorFactory generatorFactory;

    @Mock
    private Generator generator;

    @InjectMocks
    private GeneratorServiceImpl generatorService;

    @Test
    void generatePatch_returnsResult() {
        var version = RandomizerVersion.of("1.0.0-blitz-1.0");
        var preset = "Triforce Blitz";
        var randomizer = new Randomizer(version, Map.of(preset, new Preset()));
        var seed = "Test";

        var expected = new GeneratorOutput<>(new Patch(new byte[0]), new Object());
        when(generatorFactory.create(any(), any())).thenReturn(generator);
        when(randomizerRepository.load(any())).thenReturn(Optional.of(randomizer));
        when(generator.generatePatch(seed)).thenReturn(expected);

        var result = generatorService.generatePatch(version, preset, seed);

        verify(generator, atLeastOnce()).generatePatch(seed);
        assertThat(result).isEqualTo(expected);
    }


    @Test
    void generatePatch_ifRandomizerDoesNotExist_throwsException() {
        var version = RandomizerVersion.of("1.0.0-blitz-1.0");
        var preset = "Triforce Blitz";
        var seed = "Test";

        when(randomizerRepository.load(any(RandomizerVersion.class))).thenReturn(Optional.empty());

        assertThatThrownBy(() -> generatorService.generatePatch(version, preset, seed))
                .isInstanceOf(RuntimeException.class)
                .hasMessageContaining("not found");

        verify(randomizerRepository, atLeastOnce()).load(version);
    }

    @Test
    void generatePatch_ifPresetDoesNotExist_throwsException() {
        var version = RandomizerVersion.of("1.0.0-blitz-1.0");
        var randomizer = new Randomizer(version);
        var preset = "Triforce Blitz";
        var seed = "Test";

        when(randomizerRepository.load(any(RandomizerVersion.class))).thenReturn(Optional.of(randomizer));

        assertThatThrownBy(() -> generatorService.generatePatch(version, preset, seed))
                .isInstanceOf(RuntimeException.class)
                .hasMessageContaining("does not have preset");
    }
}