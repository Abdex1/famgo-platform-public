import { PartialType } from '@nestjs/swagger';
import { CreateExampleDto } from './create-example.dto';

/**
 * Update Example DTO
 * Request payload for updating examples
 * Inherits from CreateExampleDto with all fields optional
 */
export class UpdateExampleDto extends PartialType(CreateExampleDto) {}
