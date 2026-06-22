import { Test, TestingModule } from '@nestjs/testing';
import { getRepositoryToken } from '@nestjs/typeorm';
import { ExampleService } from '../example.service';
import { Example } from '../entities/example.entity';
import { NotFoundException } from '@nestjs/common';

describe('ExampleService', () => {
  let service: ExampleService;
  let mockRepository;

  const mockExample: Example = {
    id: '1',
    name: 'Test Example',
    description: 'Test Description',
    isActive: true,
    createdAt: new Date(),
    updatedAt: new Date(),
  };

  beforeEach(async () => {
    mockRepository = {
      find: jest.fn().mockResolvedValue([mockExample]),
      findOneBy: jest.fn().mockResolvedValue(mockExample),
      create: jest.fn().mockReturnValue(mockExample),
      save: jest.fn().mockResolvedValue(mockExample),
      remove: jest.fn().mockResolvedValue(null),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ExampleService,
        {
          provide: getRepositoryToken(Example),
          useValue: mockRepository,
        },
      ],
    }).compile();

    service = module.get<ExampleService>(ExampleService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('findAll', () => {
    it('should return an array of examples', async () => {
      const result = await service.findAll();
      expect(result).toEqual([mockExample]);
      expect(mockRepository.find).toHaveBeenCalled();
    });
  });

  describe('findOne', () => {
    it('should return a single example', async () => {
      const result = await service.findOne('1');
      expect(result).toEqual(mockExample);
    });

    it('should throw NotFoundException if not found', async () => {
      mockRepository.findOneBy.mockResolvedValue(null);
      await expect(service.findOne('999')).rejects.toThrow(
        NotFoundException,
      );
    });
  });

  describe('create', () => {
    it('should create and return an example', async () => {
      const createDto = { name: 'New Example' };
      const result = await service.create(createDto);
      expect(result).toEqual(mockExample);
      expect(mockRepository.create).toHaveBeenCalledWith(createDto);
    });
  });

  describe('update', () => {
    it('should update and return an example', async () => {
      const updateDto = { name: 'Updated Example' };
      const result = await service.update('1', updateDto);
      expect(result).toEqual(mockExample);
    });
  });

  describe('remove', () => {
    it('should remove an example', async () => {
      await service.remove('1');
      expect(mockRepository.remove).toHaveBeenCalledWith(mockExample);
    });
  });
});
