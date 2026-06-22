import {
  ExceptionFilter,
  Catch,
  ArgumentsHost,
  HttpException,
  HttpStatus,
  Logger,
} from '@nestjs/common';
import { Request, Response } from 'express';

/**
 * Global HTTP Exception Filter
 * Standardizes error responses across all services
 * FROM: FamGo Platform Standards
 */
@Catch(HttpException)
export class HttpExceptionFilter implements ExceptionFilter {
  private readonly logger = new Logger(HttpExceptionFilter.name);

  catch(exception: HttpException, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();
    const status = exception.getStatus();
    const exceptionResponse = exception.getResponse();

    const errorResponse = {
      statusCode: status,
      timestamp: new Date().toISOString(),
      path: request.url,
      method: request.method,
      message:
        typeof exceptionResponse === 'object' &&
        'message' in exceptionResponse
          ? exceptionResponse['message']
          : exception.message,
      error:
        typeof exceptionResponse === 'object'
          ? exceptionResponse['error']
          : 'Internal Server Error',
    };

    // Log errors
    if (status >= 500) {
      this.logger.error(
        `[${request.method}] ${request.url}`,
        exception.stack,
      );
    } else if (status >= 400) {
      this.logger.warn(`[${request.method}] ${request.url} - ${status}`);
    }

    response.status(status).json(errorResponse);
  }
}

/**
 * Catch-all exception filter for non-HTTP exceptions
 */
@Catch()
export class AllExceptionsFilter implements ExceptionFilter {
  private readonly logger = new Logger(AllExceptionsFilter.name);

  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();

    const status = HttpStatus.INTERNAL_SERVER_ERROR;

    this.logger.error(
      `[${request.method}] ${request.url}`,
      exception instanceof Error ? exception.stack : JSON.stringify(exception),
    );

    response.status(status).json({
      statusCode: status,
      timestamp: new Date().toISOString(),
      path: request.url,
      message: 'Internal server error',
      error: 'Internal Server Error',
    });
  }
}
